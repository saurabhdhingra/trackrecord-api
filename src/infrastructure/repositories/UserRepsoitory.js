import UserModel, { findById as _findById, findByIdAndUpdate } from '../models/UserModel';

class UserRepository {
    async findById(userId){
        return _findById(userId).lean();
    }

    async create(userData){
        const user = new UserModel(userData);
        return user.save();
    }

    async updateProfile(userId, data){
        return findByIdAndUpdate(userId, data, { new : true }).lean();
    }

    async addMeasurement(userId, measurementData) {
        const result = await findByIdAndUpdate(
            userId,
            { $push: {measurements: measurementData } },
            { new : true, runValidators: true }
        ).lean();
        return result.measurements.pop();
    }
}

export default UserRepository;