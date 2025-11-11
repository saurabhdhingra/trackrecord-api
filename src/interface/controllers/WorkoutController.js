class WorkoutController {
    constructor(workoutUseCases) {
        this.workoutUseCases = workoutUseCases;
    }

    // POST /workouts
    async logWorkout(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const { exercises } = req.body; // exercises is an array of { name, muscleGroup, sets: [{ reps, weight }] }

            if (!userId) {
                return res.status(401).send({ message: 'User ID is missing. Authentication required.' });
            }

            if (!exercises || exercises.length === 0) {
                return res.status(400).json({ message: 'Workout must include exercises.' });
            }

            const newWorkout = await this.workoutUseCases.logWorkout(userId, exercises);
            res.status(201).json({ message: 'Workout logged successfully.', workout: newWorkout });
        } catch (error) {
            if (error.message.includes('already logged')) {
                return res.status(409).json({ message: error.message });
            }
            res.status(400).json({ message: 'Error logging workout. Check input data.', error: error.message });
        }
    }

    // GET /workouts/report/:exerciseName
    async getEffortReport(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const exerciseName = req.params.exerciseName;

            if (!userId) {
                return res.status(401).send({ message: 'User ID is missing. Authentication required.' });
            }

            if (!exerciseName) {
                return res.status(400).json({ message: 'Exercise name is required for the report.' });
            }

            const report = await this.workoutUseCases.getPhysicalEffortReport(userId, exerciseName);
            
            res.json({
                message: `Physical effort history for ${exerciseName}.`,
                exercise: exerciseName,
                data: report 
            });
        } catch (error) {
            res.status(500).json({ message: 'Error generating report.', error: error.message });
        }
    }
}

module.exports = WorkoutController;