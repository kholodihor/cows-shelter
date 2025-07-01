provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "Terraform"
    }
  }
}

# Create an S3 bucket for file uploads
resource "aws_s3_bucket" "cows_shelter_uploads" {
  bucket = var.s3_bucket_name

  tags = {
    Name        = "Cows Shelter Uploads"
    Environment = var.environment
  }
}

# Set bucket ACL to private
# S3 bucket ownership controls - required when ACLs are disabled
resource "aws_s3_bucket_ownership_controls" "cows_shelter_uploads_ownership" {
  bucket = aws_s3_bucket.cows_shelter_uploads.id
  
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

# Enable versioning
resource "aws_s3_bucket_versioning" "cows_shelter_uploads_versioning" {
  bucket = aws_s3_bucket.cows_shelter_uploads.id
  versioning_configuration {
    status = "Enabled"
  }
}

# Enable server-side encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "cows_shelter_uploads_encryption" {
  bucket = aws_s3_bucket.cows_shelter_uploads.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Database resources are now managed by the backend module

# VPC Configuration
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames  = true
  
  tags = {
    Name = "cows-shelter-vpc"
  }
}

# Subnets
resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = "${var.aws_region}${count.index % 2 == 0 ? "a" : "b"}"
  
  tags = {
    Name = "cows-shelter-public-${count.index + 1}"
  }
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 3}.0/24"
  availability_zone = "${var.aws_region}${count.index % 2 == 0 ? "a" : "b"}"
  
  tags = {
    Name = "cows-shelter-private-${count.index + 1}"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  
  tags = {
    Name = "cows-shelter-igw"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
  
  tags = {
    Name = "cows-shelter-public-rt"
  }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Include the backend module
module "backend" {
  source = "./backend"
  
  aws_region   = var.aws_region
  environment  = var.environment
  project_name = var.project_name
  
  # MinIO configuration
  minio_endpoint   = "minio.${var.project_name}.com:9000"
  minio_access_key = "minioadmin"
  minio_secret_key = "minioadmin"
  
  # Pass other necessary variables to the backend module
  # db_username  = var.db_username
  # db_password  = var.db_password
  # db_name      = var.db_name
}

# Include the frontend module
module "frontend" {
  source = "./frontend"
  
  aws_region          = var.aws_region
  environment         = var.environment
  project_name        = var.project_name
  frontend_bucket_name = "${var.project_name}-frontend-${var.environment}"
  
  # If you have a domain name, uncomment and set these
  # domain_name       = "your-domain.com"
  # route53_zone_id   = "your-route53-zone-id"
}

# Outputs
output "s3_bucket_name" {
  value = aws_s3_bucket.cows_shelter_uploads.id
}

output "backend_url" {
  value = module.backend.alb_dns_name
  description = "URL to access the backend API"
}
