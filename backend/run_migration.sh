#!/bin/bash

# Build the migration binary
echo "Building migration binary..."
go build -o /app/migrate /app/migrations/migrate.go

# Run the migration inside the container
echo "Running migrations..."
/app/migrate
