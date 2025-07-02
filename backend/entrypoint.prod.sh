#!/bin/sh
set -e

# Wait for the database to be ready
echo "Waiting for database to be ready..."
until nc -z $DB_HOST $DB_PORT; do
  sleep 1
done
echo "Database is ready!"

# Run database migrations
echo "Running database migrations..."
/app/app migrate

# Start the application
echo "Starting application..."
/app/app

exec "$@"
