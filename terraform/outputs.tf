output "rds_connection_details" {
  description = "RDS connection details"
  value = {
    endpoint = module.backend.db_endpoint
    db_name  = module.backend.db_name
    username = module.backend.db_username
  }
  sensitive = true
}

output "s3_bucket_details" {
  description = "S3 bucket details"
  value = {
    bucket_name = aws_s3_bucket.cows_shelter_uploads.id
    region      = aws_s3_bucket.cows_shelter_uploads.region
    arn         = aws_s3_bucket.cows_shelter_uploads.arn
  }
}

output "alb_details" {
  description = "ALB details"
  value = {
    dns_name = module.backend.alb_dns_name
    zone_id  = module.backend.alb_zone_id
  }
}

output "ecr_details" {
  description = "ECR repository details"
  value = {
    repository_url = module.backend.ecr_repository_url
    repository_arn = module.backend.ecr_repository_arn
  }
}

output "ecs_details" {
  description = "ECS service details"
  value = {
    cluster_name = module.backend.ecs_cluster_name
    service_name = module.backend.ecs_service_name
  }
}

output "frontend_details" {
  description = "Frontend deployment details"
  value = {
    bucket_name = module.frontend.frontend_bucket_name
    website_endpoint = module.frontend.frontend_bucket_website_endpoint
    cloudfront_domain = module.frontend.cloudfront_domain_name
    frontend_url = module.frontend.frontend_url
  }
}

output "frontend_url" {
  description = "URL to access the frontend application"
  value       = module.frontend.frontend_url
}

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution"
  value       = module.frontend.cloudfront_distribution_id
}
