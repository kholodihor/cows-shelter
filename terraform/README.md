# Cows Shelter AWS Infrastructure

This directory contains Terraform configurations for provisioning AWS resources for the Cows Shelter application.

## Prerequisites

1. Install [Terraform](https://www.terraform.io/downloads.html) (v1.0.0 or later)
2. Configure AWS credentials with appropriate permissions
3. Install the AWS CLI and configure it with `aws configure`

## Setup

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Create a `terraform.tfvars` file with your configuration:
   ```hcl
   aws_region     = "eu-central-1"
   environment    = "dev"
   db_username    = "cows_admin"
   db_password    = "your_secure_password"
   db_name        = "cows_shelter"
   db_instance_class = "db.t3.micro"
   ```

3. Review the execution plan:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```
   
   Confirm with `yes` when prompted.

## Outputs

After applying the configuration, Terraform will output the following:

- RDS endpoint and connection details
- S3 bucket name and ARN
- VPC and security group information

## Environment Variables

Update your `.env` file with the generated values:

```
# Database Configuration
DB_HOST=<rds_endpoint>
DB_PORT=5432
DB_USER=<db_username>
DB_PASSWORD=<db_password>
DB_NAME=<db_name>

# AWS Configuration
AWS_REGION=<aws_region>
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
S3_BUCKET_NAME=<s3_bucket_name>
```

## Destroying Resources

To destroy all created resources:

```bash
terraform destroy
```

## Important Notes

- The database password is stored in the Terraform state file. Use a secure backend like S3 with state encryption.
- In production, consider using AWS Secrets Manager for database credentials.
- The RDS instance is configured with skip_final_snapshot=true for development. Change this for production.
