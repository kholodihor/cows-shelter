# Single stage build
FROM golang:1.24.2-bookworm

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    gcc \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Update go.mod and go.sum and build the application
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app .

# Make scripts executable
RUN chmod +x /app/entrypoint.sh && \
    chmod +x /app/migrate.sh && \
    ln -sf /app/migrate.sh /app/migrate

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    dos2unix \
    && dos2unix /app/entrypoint.sh \
    && chmod +x /app/entrypoint.sh \
    && chmod +x /usr/local/bin/app \
    && rm -rf /var/lib/apt/lists/*

# Expose port 8080
EXPOSE 8080

# Set the working directory
WORKDIR /app

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# Command to run the application
CMD ["app"]
