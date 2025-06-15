#!/bin/sh

# Don't exit on errors - we want to try to start the app even if migrations fail
set +e

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until nc -z postgres 5432; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Change to app directory
cd /app

# Skip migrations for now - we'll handle them separately if needed
echo "Skipping automatic migrations for now"

# Start the application
echo "Starting the application..."
exec /usr/local/bin/app "$@"
