output "minio_alb_dns_name" {
  description = "DNS name of the MinIO load balancer"
  value       = aws_lb.minio.dns_name
}

output "minio_console_url" {
  description = "URL for the MinIO console"
  value       = "http://${aws_lb.minio.dns_name}:9001"
}

output "efs_id" {
  description = "ID of the EFS filesystem for MinIO data"
  value       = aws_efs_file_system.minio_data.id
}

output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "private_subnets" {
  description = "List of private subnet IDs"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of public subnet IDs"
  value       = module.vpc.public_subnets
}
