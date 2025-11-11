class UserController {
    constructor(userUseCases) {
        this.userUseCases = userUseCases;
    }

    // POST /users/register
    async registerUser(req, res) {
        try {
            const userId = req.headers['x-user-id']; 
            const { email, birthDate } = req.body;

            if (!userId) {
                return res.status(401).send({ message: 'User ID is missing. Authentication required.' });
            }

            const newUser = await this.userUseCases.createUser(userId, email, birthDate);
            res.status(201).json({ message: 'User profile and exercise list created successfully.', user: newUser });
        } catch (error) {
            if (error.message.includes('already exists')) {
                return res.status(409).json({ message: error.message });
            }
            res.status(500).json({ message: 'Error registering user.', error: error.message });
        }
    }

    // GET /users/profile
    async getProfile(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const user = await this.userUseCases.getProfile(userId);
            if (!user) {
                return res.status(404).json({ message: 'User not found.' });
            }
            res.json(user);
        } catch (error) {
            res.status(500).json({ message: 'Error fetching profile.', error: error.message });
        }
    }

    // POST /users/measurements
    async addMeasurement(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const measurement = req.body;
            
            const newMeasurement = await this.userUseCases.addMeasurement(userId, measurement);
            res.status(201).json({ message: 'Measurement recorded successfully.', measurement: newMeasurement });
        } catch (error) {
            res.status(400).json({ message: 'Invalid data or user not found.', error: error.message });
        }
    }

    // GET /users/exercises
    async getExercises(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const exercises = await this.userUseCases.getExerciseList(userId);
            res.json(exercises);
        } catch (error) {
            res.status(500).json({ message: 'Error fetching exercise list.', error: error.message });
        }
    }
    
    // POST /users/exercises
    async addExercise(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const { name, muscleGroup, isWeighted = true } = req.body;

            if (!name || !muscleGroup) {
                return res.status(400).json({ message: 'Name and muscleGroup are required.' });
            }

            const updatedList = await this.userUseCases.addExercise(userId, { name, muscleGroup, isWeighted });
            res.status(201).json({ message: 'Exercise added.', list: updatedList.exercises });
        } catch (error) {
            res.status(500).json({ message: 'Error adding exercise.', error: error.message });
        }
    }

    // DELETE /users/exercises/:name
    async removeExercise(req, res) {
        try {
            const userId = req.headers['x-user-id'];
            const exerciseName = req.params.name;
            
            const updatedList = await this.userUseCases.removeExercise(userId, exerciseName);
            res.json({ message: `Exercise "${exerciseName}" removed.`, list: updatedList.exercises });
        } catch (error) {
            res.status(500).json({ message: 'Error removing exercise.', error: error.message });
        }
    }
}

export default UserController;
