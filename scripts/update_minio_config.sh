#!/bin/bash

# This script updates the application configuration to use the new MinIO deployment

# Check if required environment variables are set
if [ -z "$MINIO_ENDPOINT" ] || [ -z "$MINIO_ACCESS_KEY" ] || [ -z "$MINIO_SECRET_KEY" ]; then
  echo "Error: Required environment variables not set"
  echo "Please set MINIO_ENDPOINT, MINIO_ACCESS_KEY, and MINIO_SECRET_KEY"
  exit 1
fi

# Update the docker-compose.yml file
sed -i '' -e "s|MINIO_ENDPOINT=.*|MINIO_ENDPOINT=$MINIO_ENDPOINT|" \
         -e "s|MINIO_ACCESS_KEY=.*|MINIO_ACCESS_KEY=$MINIO_ACCESS_KEY|" \
         -e "s|MINIO_SECRET_KEY=.*|MINIO_SECRET_KEY=$MINIO_SECRET_KEY|" \
         -e "s|MINIO_USE_SSL=.*|MINIO_USE_SSL=true|" \
         docker-compose.yml

# Create or update .env file
echo "MINIO_ENDPOINT=$MINIO_ENDPOINT" > .env
echo "MINIO_ACCESS_KEY=$MINIO_ACCESS_KEY" >> .env
echo "MINIO_SECRET_KEY=$MINIO_SECRET_KEY" >> .env
echo "MINIO_USE_SSL=true" >> .env
echo "MINIO_BUCKET=cows-shelter" >> .env

echo "Configuration updated successfully!"
echo "Please review the changes and restart your application."
