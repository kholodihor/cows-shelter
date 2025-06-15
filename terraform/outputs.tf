output "rds_connection_details" {
  description = "RDS connection details"
  value = {
    endpoint = aws_db_instance.cows_shelter_db.endpoint
    port     = aws_db_instance.cows_shelter_db.port
    name     = aws_db_instance.cows_shelter_db.name
    username = aws_db_instance.cows_shelter_db.username
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

output "vpc_details" {
  description = "VPC details"
  value = {
    vpc_id     = aws_vpc.main.id
    cidr_block = aws_vpc.main.cidr_block
  }
}

output "security_group_id" {
  description = "RDS Security Group ID"
  value       = aws_security_group.rds_sg.id
}
