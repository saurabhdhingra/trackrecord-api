# TrackRecord API

A RESTful API for a workout tracking application built with Go.

## Overview

TrackRecord API allows users to:

- Sign up and authenticate with JWT tokens
- Browse exercises by category and muscle group
- Create and manage workout plans
- Log completed workouts
- Generate progress reports

## Requirements

- Go (1.16+)
- PostgreSQL
- [migrate](https://github.com/golang-migrate/migrate) for database migrations

## Setup

1. Clone the repository
2. Set up a PostgreSQL database
3. Update the database connection string in `cmd/api/main.go` or pass it as a flag
4. Run migrations:
   ```
   migrate -path ./migrations -database "postgres://username:password@localhost/trackrecord?sslmode=disable" up
   ```
5. Build and run the application:
   ```
   go build -o trackrecord
   ./trackrecord -port=8080 -env=development -db-dsn="postgres://username:password@localhost/trackrecord?sslmode=disable"
   ```

## API Endpoints

### Auth

- `POST /v1/users` - Register a new user
- `POST /v1/auth/login` - Login and receive JWT token

### Exercises

- `GET /v1/exercises` - List exercises (query params: category, muscle_group, page, page_size, sort)
- `GET /v1/exercises/:id` - Get a specific exercise

### Workouts

- `GET /v1/workouts` - List user's workouts (query params: page, page_size, sort)
- `POST /v1/workouts` - Create a new workout
- `GET /v1/workouts/:id` - Get a specific workout
- `PATCH /v1/workouts/:id` - Update a workout
- `DELETE /v1/workouts/:id` - Delete a workout

### Workout Logs

- `POST /v1/workout-logs` - Log a completed workout
- `GET /v1/workout-logs` - List workout logs (query params: page, page_size, sort)
- `GET /v1/workout-logs/:id` - Get a specific workout log

### Reports

- `GET /v1/reports/progress` - Generate progress report (query params: start_date, end_date, exercise_id)

## Authentication

Most endpoints require authentication. Include a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Request & Response Examples

### Register User

Request:

```json
POST /v1/users
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2023-09-01T14:30:00Z",
    "updated_at": "2023-09-01T14:30:00Z"
  }
}
```

### Create Workout

Request:

```json
POST /v1/workouts
{
  "name": "Upper Body Strength",
  "description": "Focus on chest, back, and arms",
  "schedule": "2023-09-10T10:00:00Z",
  "items": [
    {
      "exercise_id": 1,
      "sets": 3,
      "reps": 10,
      "weight": 60.5
    },
    {
      "exercise_id": 6,
      "sets": 3,
      "reps": 8,
      "weight": 0
    }
  ]
}
```

Response:

```json
{
  "workout": {
    "id": 1,
    "name": "Upper Body Strength",
    "description": "Focus on chest, back, and arms",
    "user_id": 1,
    "schedule": "2023-09-10T10:00:00Z",
    "items": [
      {
        "id": 1,
        "workout_id": 1,
        "exercise_id": 1,
        "sets": 3,
        "reps": 10,
        "weight": 60.5,
        "exercise": {
          "id": 1,
          "name": "Bench Press",
          "description": "Lie on a flat bench with a barbell and press the weight up from your chest.",
          "category": "strength",
          "muscle_group": "chest"
        }
      },
      {
        "id": 2,
        "workout_id": 1,
        "exercise_id": 6,
        "sets": 3,
        "reps": 8,
        "weight": 0,
        "exercise": {
          "id": 6,
          "name": "Pull-Ups",
          "description": "Hanging from a bar and pulling yourself up until your chin is over the bar.",
          "category": "strength",
          "muscle_group": "back"
        }
      }
    ],
    "created_at": "2023-09-01T15:00:00Z",
    "updated_at": "2023-09-01T15:00:00Z"
  }
}
```

## Project Structure

- `/cmd/api` - Application entry point and handlers
- `/internal/data` - Data models and database interaction
- `/internal/validator` - Input validation
- `/internal/auth` - Authentication utilities
- `/migrations` - Database migrations
