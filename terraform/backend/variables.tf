variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "eu-central-1"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "cows-shelter"
}

variable "environment" {
  description = "Environment (e.g., dev, staging, prod)"
  type        = string
  default     = "prod"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "List of public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "private_subnet_cidrs" {
  description = "List of private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

# Database Configuration
variable "db_instance_class" {
  description = "The instance class of the RDS instance"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "The allocated storage in GB"
  type        = number
  default     = 20
}

variable "db_username" {
  description = "The master username for the database"
  type        = string
  default     = "postgres"
}

variable "db_name" {
  description = "The name of the database"
  type        = string
  default     = "cows_shelter"
}

# MinIO Configuration
variable "minio_endpoint" {
  description = "The endpoint URL for MinIO"
  type        = string
}

variable "minio_access_key" {
  description = "The access key for MinIO"
  type        = string
  sensitive   = true
}

variable "minio_secret_key" {
  description = "The secret key for MinIO"
  type        = string
  sensitive   = true
}

# ECS Configuration
variable "desired_count" {
  description = "Number of ECS tasks to run"
  type        = number
  default     = 2
}

# Tags
variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
