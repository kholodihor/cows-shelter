variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, production)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "cows-shelter"
}

variable "frontend_bucket_name" {
  description = "S3 bucket name for frontend static files"
  type        = string
  default     = "cows-shelter-frontend"
}

variable "domain_name" {
  description = "Domain name for the frontend application (optional)"
  type        = string
  default     = ""
}

variable "route53_zone_id" {
  description = "Route53 hosted zone ID (required if domain_name is provided)"
  type        = string
  default     = ""
}
