# Cows Shelter Application

A full-stack web application for managing a cow shelter, built with React/Vite frontend and Go backend, deployed on AWS with Terraform.

## Project Structure

- `/backend` - Go backend service (REST API)
- `/frontend` - React/Vite frontend application
- `/terraform` - Infrastructure as Code (AWS resources)

## Prerequisites

- Docker and Docker Compose (for local development)
- Node.js 16+ and npm/yarn (for frontend development)
- Go 1.21+ (for backend development)
- AWS Account with appropriate permissions
- Terraform 1.0+ (for infrastructure deployment)

## Local Development

### Backend Setup

1. Copy the example environment file:
   ```bash
   cd backend
   cp .env.example .env
   ```

2. Update the `.env` file with your local configuration.

3. Start the backend service:
   ```bash
   docker-compose up -d
   ```

### Frontend Setup

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. Copy the example environment file:
   ```bash
   cp .env.example .env.local
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

## Environment Variables

### Backend (`.env`)

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_NAME=cows_shelter

# AWS Configuration
AWS_REGION=eu-north-1
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
S3_BUCKET_NAME=cows-shelter-uploads

# JWT Secret (generate a secure secret)
JWT_SECRET=your_jwt_secret_here

# App Configuration
PORT=8080
ENV=development
```

### Frontend (`.env.local`)

```
VITE_API_URL=http://localhost:8080/api
VITE_ENV=development
```

## Deployment

### Infrastructure

1. Navigate to the Terraform directory:
   ```bash
   cd terraform
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Review the execution plan:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```

### Backend Deployment

1. Build and push the Docker image:
   ```bash
   cd backend
   docker build -t your-ecr-repo/cows-shelter-backend:latest .
   docker push your-ecr-repo/cows-shelter-backend:latest
   ```

2. Update the ECS service with the new task definition.

### Frontend Deployment

1. Build the production bundle:
   ```bash
   cd frontend
   npm run build
   ```

2. Sync with S3:
   ```bash
   aws s3 sync dist/ s3://your-s3-bucket-name --delete
   ```

3. Invalidate CloudFront cache:
   ```bash
   aws cloudfront create-invalidation --distribution-id YOUR_DISTRIBUTION_ID --paths "/*"
   ```

## Contributing

1. Create a new branch for your feature/fix
2. Commit your changes with a descriptive message
3. Push to the branch
4. Create a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
