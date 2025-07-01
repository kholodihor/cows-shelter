# Database Outputs
output "db_endpoint" {
  description = "The connection endpoint for the database"
  value       = aws_db_instance.main.endpoint
}

output "db_name" {
  description = "The database name"
  value       = aws_db_instance.main.db_name
}

output "db_username" {
  description = "The master username for the database"
  value       = aws_db_instance.main.username
  sensitive   = true
}

# ALB Outputs
output "alb_dns_name" {
  description = "The DNS name of the load balancer"
  value       = aws_lb.app.dns_name
}

output "alb_zone_id" {
  description = "The zone ID of the load balancer"
  value       = aws_lb.app.zone_id
}

# ECR Outputs
output "ecr_repository_url" {
  description = "The URL of the ECR repository"
  value       = aws_ecr_repository.app.repository_url
}

output "ecr_repository_arn" {
  description = "The ARN of the ECR repository"
  value       = aws_ecr_repository.app.arn
}

# ECS Outputs
output "ecs_cluster_name" {
  description = "The name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  description = "The name of the ECS service"
  value       = aws_ecs_service.app.name
}

# VPC Outputs
output "vpc_id" {
  description = "The ID of the VPC"
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

# Security Group Outputs
output "alb_security_group_id" {
  description = "The ID of the ALB security group"
  value       = aws_security_group.alb.id
}

output "ecs_security_group_id" {
  description = "The ID of the ECS security group"
  value       = aws_security_group.ecs.id
}

output "rds_security_group_id" {
  description = "The ID of the RDS security group"
  value       = aws_security_group.rds.id
}
