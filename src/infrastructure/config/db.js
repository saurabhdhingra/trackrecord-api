import { connect } from 'mongoose';

const connectDB = async () => {
    try {
        await connect(process.env.MONGO_URI);
        console.log('MongoDB connected successfully.');
    } catch(err){
        console.log('MongoDB connection error: ', err.message);
        process.exit(1);
    }
}

export default connectDB;