import WorkoutModel from '../models/WorkoutModel';
import ExerciseListModel from '../models/ExerciseListModel';

class WorkoutRepository {
    // Exercise List methods
    async createExerciseList(userId, initialExercises) {
        const exerciseList = new ExerciseListModel({
            _id: userId,
            exercises: initialExercises,
        });
        
        return exerciseList.save();
    }

    async getExerciseList(userId) {
        return ExerciseListModel.findById(userId).lean();
    }

    async addExercise(userId, exercise){
        return ExerciseListModel.findByIdAndUpdate(
            userId,
            { $addToSet: { exercises: exercise} },
            { new: true }
        ).lean();
    }

    async removeExercise(userId, exerciseName){
        return ExerciseListModel.findByIdAndUpdate(
            userId,
            { $pull: { exercies: { name: exerciseName } } },
            { new: true }
        ).lean();
    }

    async removeExercise(userId, exerciseName) {
        return ExerciseListModel.findByIdAndUpdate(
            userId,
            { $pull: { exercises: { name: exerciseName } } }, 
            { new: true}
        ).lean();
    }


    // Workout methods

    async findWorkoutByDate(userId, date){
        const startOfDay = new Date(date);
        startOfDay.setHours(0, 0, 0, 0);
        const endOfDay = new Date(date);
        endOfDay.setHours(23, 59, 59, 999);
        
        return WorkoutModel.findOne({
            userId: userId,
            date: { $gte: startOfDay, $lte: endOfDay }
        }).lean();
    }

    async logWorkout(workoutDate) {
        const workout = new WorkoutModel(workoutDate);
        return workout.save();
    }

    async getEffortHistory(userId, exerciseName, limit = 10) {
        return WorkoutModel.find({
            userId: userId,
            'exercises.name': exerciseName
        })
        .sort({ date : -1})
        .limit(limit)
        .lean();
    }
}

module.exports = WorkoutRepository;