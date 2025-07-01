# MinIO on AWS ECS Fargate

This Terraform configuration sets up a highly available MinIO deployment on AWS ECS Fargate with the following components:

- **VPC** with public and private subnets across multiple availability zones
- **ECS Fargate** for running MinIO containers
- **EFS** for persistent storage
- **Application Load Balancer** for routing traffic
- **CloudWatch Logs** for logging
- **IAM** roles and policies with least privilege

## Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) 1.0.0 or later
2. AWS CLI configured with appropriate credentials
3. AWS account with necessary permissions

## Deployment Steps

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
   terraform apply -var="minio_root_user=your_username" -var="minio_root_password=your_secure_password"
   ```
   
   Or use a variables file:
   ```bash
   terraform apply -var-file=terraform.tfvars
   ```

4. **Access MinIO**
   - MinIO Console: http://<ALB_DNS_NAME>:9001
   - MinIO API: http://<ALB_DNS_NAME>:9000

   You can find the ALB DNS name in the Terraform outputs after deployment.

## Configuration

### Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| aws_region | AWS region to deploy resources | string | "eu-central-1" |
| project_name | Name of the project | string | "cows-shelter" |
| environment | Environment (e.g., dev, staging, prod) | string | "prod" |
| vpc_cidr | CIDR block for VPC | string | "10.0.0.0/16" |
| public_subnet_cidrs | List of public subnet CIDR blocks | list(string) | ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"] |
| private_subnet_cidrs | List of private subnet CIDR blocks | list(string) | ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"] |
| minio_root_user | MinIO root username | string | - |
| minio_root_password | MinIO root password | string | - |

## Updating the Application

1. Make your changes to the Terraform configuration
2. Run `terraform plan` to see what will be changed
3. Apply the changes with `terraform apply`

## Destroying Resources

To tear down all resources created by this configuration:

```bash
terraform destroy
```

## Security Considerations

1. **Credentials**: Never commit sensitive information like `minio_root_password` to version control. Use environment variables or a secure secret management system.

2. **Network Security**: The default configuration opens the MinIO console to the public internet. In production, consider restricting access to specific IP addresses.

3. **Encryption**: EFS volumes are encrypted at rest by default.

## Monitoring

- **CloudWatch Logs**: Logs for the MinIO containers are sent to CloudWatch Logs.
- **ECS Service Metrics**: Monitor the ECS service in the AWS Management Console.

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
