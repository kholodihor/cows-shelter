# Cows Shelter Backend

This is the backend service for the Cows Shelter application, built with Go and PostgreSQL.

## Features

- RESTful API endpoints
- PostgreSQL database integration
- File uploads to S3
- JWT authentication
- Environment-based configuration

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- AWS account (for S3 storage)
- Git

## Getting Started

### Environment Setup

1. Copy the example environment file and update it with your configuration:
   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your database and AWS credentials.

### Database Setup

1. Ensure PostgreSQL is running
2. Create a database for the application
3. Update the `.env` file with your database connection details

### Running the Application

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run migrations and start the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` by default.

## AWS Deployment

The application is ready for deployment to AWS with the following infrastructure:

### AWS Infrastructure (Terraform)

- **RDS PostgreSQL** database
- **S3 bucket** for file uploads
- **VPC** with public and private subnets
- Security groups and networking

### Deployment Steps

1. **Set up AWS credentials**:
   ```bash
   aws configure
   ```

2. **Deploy infrastructure** (from the terraform directory):
   ```bash
   cd ../terraform
   terraform init
   terraform apply
   ```

3. **Update environment variables** with the output from Terraform.

4. **Deploy the application** to an EC2 instance using the deploy script:
   ```bash
   chmod +x ../scripts/deploy.sh
   ../scripts/deploy.sh
   ```

## Database Migrations

Migrations are automatically run when the application starts. To manually run migrations:

```bash
go run main.go migrate
```

## Environment Variables

See `.env.example` for all available environment variables.

## API Documentation

API documentation is available at `/swagger/index.html` when running in development mode.

## Development

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Write tests for new features

### Testing

Run tests with:

```bash
go test ./...
```

## Deployment Script Details

The deployment script (`../scripts/deploy.sh`) performs the following actions:

1. Creates an application user
2. Sets up the application directory structure
3. Installs system dependencies (PostgreSQL client, Go)
4. Configures environment variables
5. Builds the application
6. Sets up a systemd service
7. Starts and enables the service

## Monitoring

- Check service status: `sudo systemctl status cows-shelter.service`
- View logs: `sudo journalctl -u cows-shelter.service -f`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
