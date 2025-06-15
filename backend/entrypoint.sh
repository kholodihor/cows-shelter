#!/bin/sh
set -e

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until nc -z postgres 5432; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Change to app directory
cd /app

# Run database migrations
echo "Running database migrations..."
app migrate

# Start the application
echo "Starting the application..."
exec app "$@"
