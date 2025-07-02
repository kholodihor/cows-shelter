#!/bin/bash

# Navigate to the backend directory
cd "$(dirname "$0")/.."

# Add AWS SDK dependencies
echo "Adding AWS SDK dependencies..."
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/credentials
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/aws/aws-sdk-go-v2/service/s3/types

# Tidy up the dependencies
echo "Tidying up dependencies..."
go mod tidy

echo "Dependencies updated successfully!"
