#!/bin/bash
set -e

echo "Running database migrations..."
/usr/local/bin/app migrate
