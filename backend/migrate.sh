#!/bin/bash
set -e

echo "Running database migrations..."
/app/app migrate
