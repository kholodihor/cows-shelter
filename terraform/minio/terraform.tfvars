aws_region          = "eu-central-1"
project_name        = "cows-shelter"
environment         = "prod"
vpc_cidr           = "10.0.0.0/16"
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
private_subnet_cidrs = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

# You'll need to provide these values when running terraform apply
# minio_root_user     = "your_username"
# minio_root_password = "your_secure_password"
