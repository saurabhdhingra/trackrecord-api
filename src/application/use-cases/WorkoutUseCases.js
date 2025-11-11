const { Workout } = require('../../domain/workout.entity');

class WorkoutUseCases {
    constructor(workoutRepository){
        this.workoutRepository = workoutRepository;
        this.MAX_REPORT_ENTRIES = 10;
    }

    async logWorkout(userId, exercisesData) {
        const workoutDate = new Date();

        const existingWorkout = await this.workoutRepository.findWorkoutByDate(userId, workoutDate);
        if (existingWorkout) {
            throw new Error('Workout already logged for today.');
        }

        const workout = new Workout(userId, exercisesData);

        return this.workoutRepository.logWorkout(workout);
    }

    async getPhysicalEffortReport(userId, exerciseName) {
        const rawWorkouts = await this.workoutRepository.getEffortHistory(
            userId,
            exerciseName,
            this.MAX_REPORT_ENTRIES
        );

        if (rawWorkouts.length === 0) {
            return [];
        }

        const reportData = rawWorkouts.map(rawWorkout => {
            const exerciseLog = rawWorkout.exercise.find(e => e.name === exerciseName);

            const workoutDomain = new Workout(rawWorkout.userId, rawWorkout.exercise);

            const domainExerciseLog = workoutDomain.exercises.find(e => e.name === exerciseName);

            if (!domainExerciseLog){
                return null;
            }

            const totalEffort = domainExerciseLog.calculateTotalEffort();
            
            return { 
                date: rawWorkout.date.toISOString().split('T')[0],
                exercise: exerciseName,
                totalPhysicalEffort: totalEffort,
            };
        }).filter(item => item != null)
        .reverse();

        return reportData;
    }
}

module.exports = WorkoutUseCases;