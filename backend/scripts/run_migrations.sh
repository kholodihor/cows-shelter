#!/bin/bash
set -e

# Load environment variables
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Wait for PostgreSQL to be ready
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Run migrations
echo "Running database migrations..."
if [ -f /app/migrate ]; then
  /app/migrate
else
  echo "Migration binary not found. Running migrations through the app..."
  /app/app migrate
fi

echo "Migrations completed successfully!"
