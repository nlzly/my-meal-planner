# My Meal Planner

A full-stack meal planning application with a Go backend server and React frontend.

## Project Structure

The project is organized into two main parts:

- `server/`: Go backend server
  - `api/`: API handlers
  - `db/`: Database operations (in-memory storage)
  - `models/`: Data models

- `client/`: React frontend
  - `public/`: Static files
  - `src/`: React components and application logic
    - `components/`: Reusable UI components

## Features

- Weekly meal planning
- Add, view, and delete meals
- Organize meals by day and type (breakfast, lunch, dinner)
- RESTful API for meal management

## Getting Started

### Prerequisites

- Go (1.16+)
- Node.js (14+) and npm

### Running the Server

```bash
cd server
go run main.go
```

The server will start on http://localhost:8080

### Running the Client

```bash
cd client
npm install
npm start
```

The client will start on http://localhost:3000

## API Endpoints

- `GET /api/health`: Health check endpoint
- `GET /api/meals`: Get all meals
- `POST /api/meals`: Create a new meal
- `GET /api/meals/:id`: Get a meal by ID
- `PUT /api/meals/:id`: Update a meal
- `DELETE /api/meals/:id`: Delete a meal
