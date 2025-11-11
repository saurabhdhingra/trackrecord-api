import { Schema, model } from 'mongoose';

const exerciseSchema = new Schema({
    name: { type: String, required: true },
    muscleGroup: { type: String, required: true },
    isWeighted: { type: Boolean, default: true }
}, { _id: false });

const exerciseListSchema = new Schema({
    _id: { type: String, required: true, ref: 'User' },
    exercises: [exerciseSchema]
}, {
    collection: 'exercise_lists',
    timestamps: true
});

const ExerciseListModel = model('ExerciseList', exerciseListSchema);
export default ExerciseListModel;