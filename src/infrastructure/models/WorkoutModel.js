import { Schema, model } from 'mongoose';

const setSchema = new Schema({
    repetitions: { type: Number, required: true },
    weight: { type: Number, required: true, default: 0},
}, { _id: false});

const exerciseLogSchema = new Schema({
    name: { type: String, required: true },
    muscleGroup: { type: String, required: true }, 
    sets: [setSchema]
}, { _id: false });

const workoutSchema = new Schema({
    userId: { type: String, required: true, ref: 'User'}, 
    date: { type: Date, required: true, default: Date.now },
    exercises: [exerciseLogSchema]
}, {
    collection: 'workouts',
    timestamps: true
});

workoutSchema.index({ userId: 1, date: 1}, {unique: true, partialFilterExpression: { date: { $exists: true } } });

const WorkoutModel = model('Workout', workoutSchema);
export default WorkoutModel; 