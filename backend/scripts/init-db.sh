#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f "../.env" ]; then
    export $(grep -v '^#' ../.env | xargs)
else
    echo "Error: .env file not found in the parent directory."
    exit 1
fi

# Check if required environment variables are set
for var in DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME; do
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set in .env file"
        exit 1
    fi
done

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
    echo "PostgreSQL is not ready yet. Retrying in 5 seconds..."
    sleep 5
done

echo "PostgreSQL is ready!"

# Run database migrations
echo "Running database migrations..."
cd .. && go run main.go migrate

echo "Database initialization completed successfully!"
