import { Router } from 'express';
const router = Router();

export default (workoutController) => {
    // Log a new workout (one per day enforcement is in the Use Case)
    router.post('/', workoutController.logWorkout.bind(workoutController));
    
    // Report endpoint (Physical Effort History)
    router.get('/report/:exerciseName', workoutController.getEffortReport.bind(workoutController));

    return router;
};