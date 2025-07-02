package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

// InitMinio initializes the MinIO client
func InitMinio() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	// Initialize MinIO client with appropriate credentials
	var client *minio.Client
	var err error

	// Check if we're running in AWS (using S3 endpoint) or local MinIO
	isAWS := endpoint == "s3.amazonaws.com" || 
			 endpoint == "s3.eu-north-1.amazonaws.com" ||
			 endpoint == "" // If no endpoint is provided, assume AWS S3

	if isAWS {
		// Use AWS SDK's default credential chain which will automatically handle IAM roles
		log.Println("Using AWS IAM role-based authentication for S3")
		
		// Use AWS SDK to get credentials from the default chain
		awsCfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion("eu-north-1"), // Explicitly set the region
		)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %v", err)
		}
		
		// Get credentials from the AWS config
		awsCreds, err := awsCfg.Credentials.Retrieve(context.Background())
		if err != nil {
			return fmt.Errorf("failed to retrieve AWS credentials: %v", err)
		}

		// Log the AWS credentials source for debugging
		log.Printf("Using AWS credentials from: %s", awsCreds.Source)
		
		// Create MinIO client with AWS credentials
		// Use the regional endpoint for S3
		client, err = minio.New("s3.eu-north-1.amazonaws.com", &minio.Options{
			Creds:  credentials.NewStaticV4(awsCreds.AccessKeyID, awsCreds.SecretAccessKey, awsCreds.SessionToken),
			Secure: true,
			Region: "eu-north-1",
		})
	} else {
		// For local MinIO or when explicit credentials are provided
		if accessKey == "" || secretKey == "" {
			return fmt.Errorf("MINIO_ACCESS_KEY and MINIO_SECRET_KEY must be provided for non-AWS S3 endpoints")
		}
		
		log.Println("Using static credentials for MinIO")
		client, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
	}
	
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	MinioClient = client

	// Create bucket if it doesn't exist
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		return fmt.Errorf("MINIO_BUCKET environment variable is not set")
	}

	log.Printf("Checking if bucket '%s' exists in S3\n", bucketName)
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	log.Printf("Bucket '%s' exists: %v\n", bucketName, exists)

	// For AWS S3, we don't need to create the bucket as it should already exist
	// and the IAM role should have permissions to access it
	if !exists {
		log.Printf("Bucket '%s' does not exist. In AWS S3, the bucket should be created in advance.\n", bucketName)
		return fmt.Errorf("bucket %s does not exist. Please create it in the AWS S3 console", bucketName)
	}

	// Verify we can list objects in the bucket as a test
	log.Printf("Testing list objects in bucket '%s'\n", bucketName)
	objectCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive: true,
		MaxKeys:   1,
	})
	
	// Check if we can read from the channel without error
	_, ok := <-objectCh
	if !ok {
		log.Printf("Warning: No objects found in bucket or error listing objects\n")
	} else {
		log.Printf("Successfully accessed bucket '%s'\n", bucketName)
	}
	
	// Drain the channel to avoid goroutine leak
	for range objectCh {
		// Just drain the channel
	}

	log.Println("MinIO client initialized successfully")
	return nil
}
