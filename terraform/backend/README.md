# Cows Shelter Backend Infrastructure

This Terraform configuration sets up the AWS infrastructure for the Cows Shelter backend application, including:

- VPC with public and private subnets
- RDS PostgreSQL database
- ECS Fargate service
- Application Load Balancer
- Security groups and IAM roles
- CloudWatch logging

## Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) 1.0.0 or later
2. AWS CLI configured with appropriate credentials
3. AWS account with necessary permissions

## Getting Started

1. **Initialize Terraform**
   ```bash
   terraform init
   ```

2. **Review the execution plan**
   ```bash
   terraform plan
   ```

3. **Apply the configuration**
   ```bash
   terraform apply -var-file=terraform.tfvars
   ```

## Configuration

### Required Variables

- `minio_endpoint`: The endpoint URL for MinIO
- `minio_access_key`: The access key for MinIO
- `minio_secret_key`: The secret key for MinIO

### Optional Variables

See `variables.tf` for all available variables and their default values.

## Deployment

### Building and Pushing the Docker Image

1. Build the Docker image:
   ```bash
   docker build -t cows-shelter-backend -f Dockerfile.prod .
   ```

2. Tag the image for ECR:
   ```bash
   docker tag cows-shelter-backend:latest $(aws ecr describe-repositories --repository-names cows-shelter-app --query 'repositories[0].repositoryUri' --output text):latest
   ```

3. Authenticate Docker to your ECR registry:
   ```bash
   aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query 'Account' --output text).dkr.ecr.eu-central-1.amazonaws.com
   ```

4. Push the image to ECR:
   ```bash
   docker push $(aws ecr describe-repositories --repository-names cows-shelter-app --query 'repositories[0].repositoryUri' --output text):latest
   ```

### Deploying the Application

After pushing the Docker image, the ECS service will automatically start new tasks with the latest image.

## Monitoring

- **ECS Service Metrics**: Available in the AWS ECS console
- **Application Logs**: Available in CloudWatch Logs under `/ecs/cows-shelter-app`
- **Database Metrics**: Available in the AWS RDS console

## Cleanup

To destroy all resources created by this configuration:

```bash
terraform destroy -var-file=terraform.tfvars
```

## Security Considerations

- Database credentials are stored in AWS Secrets Manager
- All resources are tagged with the environment and project name
- Security groups are configured to allow only necessary traffic
- Encryption at rest is enabled for RDS and EBS volumes

## Troubleshooting

### Common Issues

1. **Permission Denied Errors**
   - Ensure the IAM user/role has the necessary permissions
   - Check the ECS task execution role policies

2. **Container Failing to Start**
   - Check the ECS service events in the AWS Console
   - Review the CloudWatch logs for the failing container

3. **Connectivity Issues**
   - Verify security group rules allow traffic on the required ports
   - Check route tables and network ACLs

## Support

For issues and feature requests, please open an issue in the repository.
