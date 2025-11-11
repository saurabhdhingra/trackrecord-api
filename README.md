# Trackrecord API

This is a robust backend application for logging and analyzing physical workouts and body measurements. It is built with Node.js, Express, and MongoDB (via Mongoose), following the Clean Architecture principles to ensure maintainability, testability, and separation of concerns.

## Architecture Overview

The project is structured into distinct layers as per the Clean Architecture pattern:

1. **Domain**: Contains core business entities ('User', 'Workout', 'Set', 'Measurement') and fundamental rules (e.g., Physical Effort calculation).

2. **Application** (Use Cases): Contains the application's specific business logic (e.g., 'logWorkout', 'getPhysicalEffortReport', 'createUser'). It orchestrates data flow between the domain and infrastructure layers.

3. **Infrastructure**: Handles all external concerns, specifically MongoDB interaction (Mongoose Models and Repositories) and connection configuration.

4. **Interface**: Manages the external interfaces, such as Express controllers and routes, handling HTTP requests and responses.

5. **Composition Root** ('server.js'): Initializes all layers and manages dependency injection.

## Features

· **Firebase UID Integration**: Uses the Firebase Authentication UID as the primary user identifier (_id) for secure data segregation.

· **Single Workout per Day**: Enforces a business rule allowing users to log only one workout entry per calendar day.

· **Custom Exercise Lists**: Users maintain a personalized, unique list of exercises they can select from.

· **Physical Effort Reporting**: Calculates a weighted physical effort metric (reps * weight) for any given exercise and provides a historical array for charting.

· **Body Measurement Tracking**: Allows logging of various body measurements (girths, height, weight, etc.).

Folder Structure : 

```
workout-tracker-api/
├── .env.example
├── package.json
├── server.js
└── src/
    ├── application/
    │   └── use-cases/
    │       ├── UserUseCases.js        (User/Exercise List Business Logic)
    │       └── WorkoutUseCases.js     (Workout/Report Business Logic)
    │
    ├── domain/
    │   ├── user.entity.js           (User and Measurement Entities)
    │   └── workout.entity.js        (Workout, ExerciseLog, and Set Entities)
    │
    ├── infrastructure/
    │   ├── config/
    │   │   └── db.js                  (MongoDB Connection)
    │   ├── models/
    │   │   ├── ExerciseListModel.js   (Mongoose Model)
    │   │   ├── UserModel.js           (Mongoose Model)
    │   │   └── WorkoutModel.js        (Mongoose Model)
    │   └── repositories/
    │       ├── UserRepository.js      (Data Access for Users)
    │       └── WorkoutRepository.js   (Data Access for Workouts/Exercises)
    │
    └── interface/
        ├── controllers/
        │   ├── UserController.js      (Handles User/Exercise HTTP requests)
        │   └── WorkoutController.js   (Handles Workout/Report HTTP requests)
        └── routes/
            ├── user.routes.js         (Express routes for /api/users)
            └── workout.routes.js      (Express routes for /api/workouts)
```

## Prerequisites

Before running the application, ensure you have the following installed:

· **Node.js** (v18 or higher recommended)

· **npm** (Node Package Manager)

· **MongoDB Atlas or Local Instance** (A running MongoDB instance)

## Installation and Setup

### 1. Clone the repository

```
git clone <repository-url>
cd workout-tracker-api
```

### 2. Install dependencies

```
npm install
```

### 3. Configure environment variables

Copy the example environment file and fill in your connection details.

```
cp ./.env.example ./.env
```

Edit the newly created .env file:

## .env file content
```
MONGO_URI=mongodb+srv://<username>:<password>@clustername.mongodb.net/<database>?retryWrites=true&w=majority
PORT=3000
```

4. Run the application

You can start the server in development mode (using nodemon) or production mode:

```
# Development (requires nodemon package)
npm run dev
```

```
# Production
npm start
```

The server will typically run on http://localhost:3000.

## API Endpoints

All API endpoints require the client to pass the authenticated Firebase User ID in a custom header: x-user-id.

**User Management** ('/api/users')

```
Method      Endpoint                    Description                         Request Body

POST        /api/users/register         Creates a new user profile and      { "email": "user@example.
                                        initializes their custom exercise   com", "birthDate": "1990-01-01" }
                                        list with default exercises.

GET         /api/users/profile          Retrieves the user's profile and    None
                                        body measurement history.

POST        /api/users/measurements     Adds a new body measurement entry   { "weight": 75, "height": 180, 
                                        for the user.                       "armsGirth": 35, ... }

GET         /api/users/exercises        Gets the user's current list        None
                                        of available exercises.

POST        /api/users/exercises        Adds a new custom exercise to the   { "name": "Cable Bicep Curls",
                                        user's list.                        "muscleGroup": "Biceps", "isWeighted": 
                                                                            true }
DELETE      /api/users/exercises/:name  Removes a custom exercise from      None (Name in URL param)
                                        the user's list.
```

**Workout Logging and Reporting** ('/api/workouts')

```
Method              Endpoint                Description                         Request Body

POST                /api/workouts           Logs a complete workout             { "exercises": [ { "name": "Bench Press",
                                            for the day. Fails if a workout     "muscleGroup": "Chest", "sets":
                                            is already logged today.            [{ "repetitions": 12, "weight": 60 }, ...] } ] }

GET                 /api/workouts/          Retrieves the physical effort       None (Exercise Name in URL param)
                    report/:exerciseName    history (Total Effort = Reps * 
                                            Weight) for a specific exercise 
                                            over the last 10 workouts.
```

## Acknowledgement
https://roadmap.sh/projects/fitness-workout-tracker