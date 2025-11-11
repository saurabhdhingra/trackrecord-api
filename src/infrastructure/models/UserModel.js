import mongoose from 'mongoose';

const measurementSchema = new mongoose.Schema({
    weight: { type : Number },
    height: { type : Number },
    leftArmGirth: { type : Number },
    rightArmGirth: { type : Number },
    leftThighGirth: { type : Number },
    rightThighGirth: { type : Number },
    waist: { type: Number },
    shoulders: { type: Number },
    chest: { type: Number },
    date: { type: Date, default: Date.now },
}, { _id : false });

const userSchema = new mongoose.Schema({
    _id : { type: String, required: true },
    email: { type: String, required: true, unique: true},
    birthDate: { type: Date },
    measurements: [measurementSchema]
}, {
    collection: 'users',
    timestamps: true
})

const UserModel = mongoose.model('User', userSchema);
module.exports = UserModel;