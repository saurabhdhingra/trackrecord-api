import { Measurement } from '../../domain/user.entity';


const INITIAL_EXERCISES = [
    { name: 'Barbell Bench Press', muscleGroup: 'Chest' },
    { name: 'Lateral Raise (Dumbbells)', muscleGroup: 'Shoulders' },
    { name: 'Bicep Curl (Barbell)', muscleGroup: 'Biceps' },
    { name: 'Squat (Barbell)', muscleGroup: 'Legs' },
    { name: 'Push-ups', muscleGroup: 'Chest', isWeighted: false }
];

class UserUseCases {
    constructor(userRepository, workoutRepository) {
        this.userRepository = userRepository;
        this.workoutRepository = workoutRepository;
    }

    // User management

    async createUser(userId, email, birthDate) {
        // 1. Check if user already exists
        const existingUser = await this.userRepository.findById(userId);
        if (existingUser) {
            throw new Error('User already exists.');
        }

        const userData = {
            _id: userId,
            email,
            birthDate
        };
        const newUser = await this.userRepository.create(userData);

        await this.workoutRepository.createExerciseList(userId, INITIAL_EXERCISES);

        return newUser;
    }

    async updateProfile(userId, data) {
        return this.userRepository.updateProfile(userId, data);
    }

    async addMeasurement(userId, measurementData) {
        const measurement = new Measurement(
            measurementData.weight,
            measurementData.height,
            measurementData.armsGirth,
            measurementData.thighGirth,
            measurementData.waist,
            measurementData.shoulders,
            measurementData.chest
        );
        return this.userRepository.addMeasurement(userId, measurement);
    }

    async getProfile(userId) {
        return this.userRepository.findById(userId);
    }

    // Exercise List Management 

    async getExerciseList(userId) {
        const list = await this.workoutRepository.getExerciseList(userId);
        return list ? list.exercises : [];
    }

    async addExercise(userId, exercise) {
        return this.workoutRepository.addExercise(userId, exercise);
    }

    async removeExercise(userId, exerciseName) {
        return this.workoutRepository.removeExercise(userId, exerciseName);
    }
}

export default UserUseCases;