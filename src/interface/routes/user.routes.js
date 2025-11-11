import { Router } from 'express';
const router = Router();

export default (userController) => {
    // Basic User Operations
    router.post('/register', userController.registerUser.bind(userController));
    router.get('/profile', userController.getProfile.bind(userController));
    
    // Body Measurements
    router.post('/measurements', userController.addMeasurement.bind(userController));

    // Custom Exercise List Management
    router.get('/exercises', userController.getExercises.bind(userController));
    router.post('/exercises', userController.addExercise.bind(userController));
    router.delete('/exercises/:name', userController.removeExercise.bind(userController));

    return router;
};