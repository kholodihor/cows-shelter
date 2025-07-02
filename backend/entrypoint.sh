#!/bin/sh

# Enable error checking after this point
set -e

# Debug information
echo "=== Debug: Entrypoint started ==="
echo "Current user: $(whoami)"
echo "Working directory: $(pwd)"
echo "Environment variables:"
printenv | sort
echo "=== Debug: File permissions in /app ==="
ls -la /app

# Function to check if database is ready
db_ready() {
    echo "Checking if database is ready..."
    nc -z -w 2 "${DB_HOST:-postgres}" "${DB_PORT:-5432}"
}

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until db_ready; do
  echo "Waiting for PostgreSQL at ${DB_HOST:-postgres}:${DB_PORT:-5432}..."
  sleep 2
done

echo "PostgreSQL is up - running migrations..."

# Change to app directory
cd /app

# Run database migrations using the Go binary directly
echo "Running database migrations..."
if [ -f "/app/app" ]; then
    if /app/app migrate; then
        echo "Migrations completed successfully"
    else
        echo "Warning: Database migrations failed. The application will start but may not function correctly."
    fi
else
    echo "Error: Application binary not found. Cannot run migrations."
    exit 1
fi

# Create admin user if not exists
echo "Ensuring admin user exists..."
if [ -n "${ADMIN_EMAIL}" ] && [ -n "${ADMIN_PASSWORD}" ]; then
    echo "Creating/updating admin user..."
    if ! /app/app create-admin -email "${ADMIN_EMAIL}" -password "${ADMIN_PASSWORD}"; then
        echo "Warning: Failed to create/update admin user. The application will start but admin access may not be available."
    fi
else
    echo "Warning: ADMIN_EMAIL and/or ADMIN_PASSWORD not set. Skipping admin user creation."
fi

# Start the application
echo "Starting the application..."
exec /app/app "$@"
