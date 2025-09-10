# Cows Shelter Application

A full-stack web application for managing a cow shelter, built with React/Vite frontend and Go backend, deployed on Fly.io with Cloudinary for image storage.

## 🚀 Live Application

- **Backend API**: Fly.io
- **Database**: Neon PostgreSQL (Serverless)
- **Image Storage**: Cloudinary CDN

## Project Structure

- `/backend` - Go backend service (REST API)
- `/frontend` - React/Vite frontend application

## Prerequisites

- Node.js 18+ and npm (for frontend development)
- Go 1.21+ (for backend development)
- Fly.io CLI (for deployment)
- Cloudinary account (for image storage)

## 🛠️ Local Development

### Backend Setup

1. Copy the example environment file:
   ```bash
   cd backend
   cp .env.example .env
   ```

2. Update the `.env` file with your configuration:
   ```bash
   # Database Configuration (Neon PostgreSQL)
   DATABASE_URL=your_neon_database_url
   
   # Cloudinary Configuration
   CLOUDINARY_CLOUD_NAME=your_cloud_name
   CLOUDINARY_API_KEY=your_api_key
   CLOUDINARY_API_SECRET=your_api_secret
   STORAGE_TYPE=cloudinary
   
   # JWT Secret (generate a secure secret)
   JWT_SECRET=your_jwt_secret_here
   
   # Admin Configuration
   ADMIN_EMAIL=admin@cows-shelter.com
   ADMIN_PASSWORD=your_secure_password
   
   # App Configuration
   PORT=8080
   ENV=development
   GIN_MODE=debug
   ```

3. Run the backend locally:
   ```bash
   go mod tidy
   go run main.go
   ```

### Frontend Setup

1. Install dependencies:
   ```bash
   cd frontend
   npm install
   ```

2. The frontend is configured to use the deployed backend at:
   ```
   VITE_API_URL=https://backend-fragrant-star-8901.fly.dev
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

## 🌐 Deployment

### Backend Deployment (Fly.io)

The backend is deployed on Fly.io with the following configuration:

1. **Database**: Neon PostgreSQL (serverless, auto-scaling)
2. **Image Storage**: Cloudinary (CDN with global distribution)
3. **Hosting**: Fly.io (global edge deployment)

To deploy updates:

```bash
cd backend
flyctl deploy --app backend-fragrant-star-8901
```

### Environment Variables (Production)

The following secrets are configured in Fly.io:

```bash
flyctl secrets set \
  DATABASE_URL="your_neon_database_url" \
  CLOUDINARY_CLOUD_NAME="your_cloud_name" \
  CLOUDINARY_API_KEY="your_api_key" \
  CLOUDINARY_API_SECRET="your_api_secret" \
  STORAGE_TYPE="cloudinary" \
  JWT_SECRET="your_jwt_secret" \
  ADMIN_EMAIL="admin@cows-shelter.com" \
  ADMIN_PASSWORD="your_secure_password" \
  --app backend-fragrant-star-8901
```

### Frontend Deployment

The frontend can be deployed to any static hosting service (Vercel, Netlify, etc.) and is configured to connect to the production backend.

## 🏗️ Architecture

- **Frontend**: React + TypeScript + Vite + Tailwind CSS
- **Backend**: Go + Gin + GORM
- **Database**: PostgreSQL (Neon - serverless)
- **File Storage**: Cloudinary (images, documents)
- **Deployment**: Fly.io (backend), Static hosting (frontend)
- **Authentication**: JWT tokens

## 📝 API Endpoints

- `GET /api/health` - Health check
- `POST /api/login` - Admin authentication
- `GET /api/news` - Get news articles
- `POST /api/news` - Create news article
- `GET /api/gallery` - Get gallery images
- `POST /api/gallery` - Upload gallery image
- `GET /api/excursions` - Get excursions
- `POST /api/excursions` - Create excursion
- `GET /api/partners` - Get partners
- `POST /api/partners` - Create partner
- `GET /api/contacts` - Get contact information
- `GET /api/reviews` - Get reviews

## 🔧 Features

- ✅ Multilingual support (Ukrainian/English)
- ✅ Admin panel for content management
- ✅ Image upload with Cloudinary integration
- ✅ Responsive design
- ✅ SEO optimized
- ✅ Global CDN for fast image delivery
- ✅ Serverless database with auto-scaling
- ✅ Production-ready deployment

## Contributing

1. Create a new branch for your feature/fix
2. Commit your changes with a descriptive message
3. Push to the branch
4. Create a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
