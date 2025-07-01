terraform {
  backend "s3" {
    bucket         = "cows-shelter-tfstate-1751359198"
    key            = "terraform.tfstate"
    region         = "eu-north-1"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }
}