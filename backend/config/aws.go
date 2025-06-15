package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kholodihor/cows-shelter-backend/utils"
)

var (
	S3Session *s3.S3
	AwsRegion = utils.GetEnv("AWS_REGION", "eu-central-1")
)

func init() {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(AwsRegion),
			Credentials: credentials.NewEnvCredentials(),
		},
	)

	if err != nil {
		panic(err)
	}

	S3Session = s3.New(sess)
}
