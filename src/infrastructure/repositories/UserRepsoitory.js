const UserModel = require('../models/UserModel');

class UserRepository {
    async findById(userId){
        return UserModel.findById(userId).lean();
    }

    async create(userData){
        const user = new UserModel(userDate);
        return user.save();
    }

    async updateProfile(userId, data){
        return UserModel.findByIdAndUpdate(userId, Date, { new : true }).lean();
    }

    async addMeasurement(userId, measurementData) {
        const result = await UserModel.findByIdAndUpdate(
            userId,
            { $push: {measurements: measurementData } },
            { new : true, runValidators: true }
        ).lean();
        return result.measurements.pop();
    }
}

module.exports = UserRepository;