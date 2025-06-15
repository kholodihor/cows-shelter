#!/bin/bash

# Exit on error
set -e

# Configuration
APP_NAME="cows-shelter"
APP_USER="cows"
APP_DIR="/opt/$APP_NAME"
SERVICE_NAME="$APP_NAME.service"
ENV_FILE="$APP_DIR/backend/.env"

# Create application user if it doesn't exist
if ! id "$APP_USER" &>/dev/null; then
    echo "Creating application user: $APP_USER"
    sudo useradd -r -s /bin/false "$APP_USER"
fi

# Create application directory
sudo mkdir -p "$APP_DIR"
sudo chown -R $USER:$USER "$APP_DIR"
chmod 755 "$APP_DIR"

# Copy application files
echo "Copying application files..."
rsync -av --exclude='.git' --exclude='node_modules' --exclude='.env' ./ "$APP_DIR/"

# Install system dependencies
echo "Installing system dependencies..."
sudo apt-get update
sudo apt-get install -y postgresql-client

# Install Go if not installed
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
    source ~/.profile
    rm go1.21.0.linux-amd64.tar.gz
fi

# Set up environment variables
echo "Setting up environment variables..."
if [ ! -f "$ENV_FILE" ]; then
    cp "$APP_DIR/backend/.env.example" "$ENV_FILE"
    echo "Please update the configuration in $ENV_FILE"
    exit 1
fi

# Build the application
echo "Building the application..."
cd "$APP_DIR/backend"
go mod download
go build -o "$APP_DIR/$APP_NAME"

# Set up systemd service
echo "Setting up systemd service..."
sudo tee "/etc/systemd/system/$SERVICE_NAME" > /dev/null <<EOL
[Unit]
Description=Cows Shelter Backend
After=network.target

[Service]
Type=simple
User=$APP_USER
WorkingDirectory=$APP_DIR/backend
EnvironmentFile=$ENV_FILE
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=$APP_NAME

[Install]
WantedBy=multi-user.target
EOL

# Set permissions
sudo chown -R $APP_USER:$APP_USER "$APP_DIR"
sudo chmod -R 750 "$APP_DIR"

# Reload systemd and start the service
echo "Starting the service..."
sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

echo "Deployment completed successfully!"
echo "Application is running as a systemd service: $SERVICE_NAME"
echo "Check status with: sudo systemctl status $SERVICE_NAME"
echo "View logs with: sudo journalctl -u $SERVICE_NAME -f"
