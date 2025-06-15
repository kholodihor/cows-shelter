provider "aws" {
  region = var.aws_region
}

# Create an S3 bucket for file uploads
resource "aws_s3_bucket" "cows_shelter_uploads" {
  bucket = var.s3_bucket_name
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

  tags = {
    Name        = "Cows Shelter Uploads"
    Environment = var.environment
  }
}

# Create an RDS PostgreSQL database
resource "aws_db_instance" "cows_shelter_db" {
  identifier             = "cows-shelter-db"
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = var.db_instance_class
  allocated_storage      = var.db_allocated_storage
  storage_type           = "gp3"
  storage_encrypted      = true
  username               = var.db_username
  password               = var.db_password
  db_name                = var.db_name
  parameter_group_name   = aws_db_parameter_group.cows_shelter_pg.name
  vpc_security_group_ids = [aws_security_group.rds_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.cows_shelter_subnet_group.name
  skip_final_snapshot    = true
  multi_az               = var.environment == "production"
  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "Mon:04:00-Mon:05:00"
  
  tags = {
    Name        = "cows-shelter-db"
    Environment = var.environment
  }
}

resource "aws_db_parameter_group" "cows_shelter_pg" {
  name   = "cows-shelter-pg"
  family = "postgres15"

  parameter {
    name  = "log_connections"
    value = "1"
  }
}

resource "aws_db_subnet_group" "cows_shelter_subnet_group" {
  name       = "cows-shelter-subnet-group"
  subnet_ids = aws_subnet.private.*.id

  tags = {
    Name = "Cows Shelter DB Subnet Group"
  }
}

# Security Group for RDS
resource "aws_security_group" "rds_sg" {
  name        = "cows-shelter-rds-sg"
  description = "Security group for RDS"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # In production, restrict this to your application servers
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "cows-shelter-rds-sg"
  }
}

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

# Outputs
output "rds_endpoint" {
  value = aws_db_instance.cows_shelter_db.endpoint
}

output "s3_bucket_name" {
  value = aws_s3_bucket.cows_shelter_uploads.id
}
