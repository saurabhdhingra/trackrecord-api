require('dotenv').config();
import express, { json } from 'express';
import connectDB from './src/infrastructure/config/db';

// Infrastructure Imports
import UserRepository from './src/infrastructure/repositories/UserRepository';
import WorkoutRepository from './src/infrastructure/repositories/WorkoutRepository';

// Application Imports (Use Cases)
import UserUseCases from './src/application/use-cases/UserUseCases';
import WorkoutUseCases from './src/application/use-cases/WorkoutUseCases';

// Interface Imports (Controllers & Routes)
import UserController from './src/interface/controllers/UserController';
import WorkoutController from './src/interface/controllers/WorkoutController';
import userRoutes from './src/interface/routes/user.routes';
import workoutRoutes from './src/interface/routes/workout.routes';

const app = express();

// Middleware
app.use(json());

// --- Dependency Injection and Composition Root ---
// 1. Initialize Repositories (Data Access)
const userRepository = new UserRepository();
const workoutRepository = new WorkoutRepository();

// 2. Initialize Use Cases (Business Logic)
const userUseCases = new UserUseCases(userRepository, workoutRepository);
const workoutUseCases = new WorkoutUseCases(workoutRepository);

// 3. Initialize Controllers (Request Handling)
const userController = new UserController(userUseCases);
const workoutController = new WorkoutController(workoutUseCases);

// 4. Register Routes
app.use('/api/users', userRoutes(userController));
app.use('/api/workouts', workoutRoutes(workoutController));

// Default Route
app.get('/', (req, res) => {
    res.send('Workout Tracker API is running.');
});

// --- Server Setup ---
const PORT = process.env.PORT || 3000;

connectDB().then(() => {
    app.listen(PORT, () => {
        console.log(`Server is running on port ${PORT}`);
        console.log('To run a sample request, set an "x-user-id" header (your Firebase UID).');
    });
});