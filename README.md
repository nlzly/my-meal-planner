# My Meal Planner

A simple meal planning application with Google OAuth authentication.

## Features

- Google OAuth login
- Weekly meal planning
- Add, view, and delete meals
- Protected API endpoints

## Setup

### Prerequisites

- Node.js and npm for the client
- Go for the server
- Google OAuth credentials

### Getting Google OAuth Credentials

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Credentials"
4. Click "Create Credentials" > "OAuth client ID"
5. Select "Web application" as the application type
6. Add authorized JavaScript origins:
   - For development: `http://localhost:3000`
   - For production: Your domain
7. Add authorized redirect URIs:
   - For development: `http://localhost:8080/auth/google/callback`
   - For production: Your domain + `/auth/google/callback`
8. Click "Create" and note your Client ID and Client Secret

### Environment Variables

#### Client

Create a `.env` file in the `client` directory with:

```
REACT_APP_GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID_HERE
```

#### Server

Set the following environment variables:

```
GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID_HERE
GOOGLE_CLIENT_SECRET=YOUR_GOOGLE_CLIENT_SECRET_HERE
OAUTH_REDIRECT_URL=http://localhost:8080/auth/google/callback
JWT_SECRET=your-jwt-secret-key
```

## Running the Application

### Server

```bash
cd server
go mod tidy
go run main.go
```

The server will start on port 8080.

### Client

```bash
cd client
npm install
npm start
```

The client will start on port 3000.

## Usage

1. Open your browser and navigate to `http://localhost:3000`
2. Click the Google Sign-In button to authenticate
3. Once authenticated, you can view, add, and delete meals in your meal planner
4. Your authentication token will be stored in localStorage and used for API requests

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/meals` - Get all meals (protected)
- `POST /api/meals` - Create a new meal (protected)
- `GET /api/meals/:id` - Get a meal by ID (protected)
- `PUT /api/meals/:id` - Update a meal (protected)
- `DELETE /api/meals/:id` - Delete a meal (protected)
- `GET /auth/google/login` - Initiate Google OAuth flow
- `POST /auth/google/callback` - Process Google OAuth callback
