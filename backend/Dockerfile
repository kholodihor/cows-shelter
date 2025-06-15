# Single stage build
FROM golang:1.22.6-alpine3.19

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache bash gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app .

# Install runtime dependencies
RUN apk add --no-cache \
    bash \
    dos2unix \
    && dos2unix /app/entrypoint.sh \
    && chmod +x /app/entrypoint.sh \
    && chmod +x /usr/local/bin/app

# Expose port 8080
EXPOSE 8080

# Set the working directory
WORKDIR /app

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# Command to run the application
CMD ["app"]
