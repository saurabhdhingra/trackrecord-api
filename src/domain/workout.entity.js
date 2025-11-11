class Set {
    constructor(repetitions, weight){
        this.repetitions = repetitions;
        this.weight = weight;
    }

    calculateEffort() {
        if (this.weight > 0){
            return this.repetitions * this.weight;
        }else {
            return this.weight;
        }
    }
}

class ExerciseLog {
    constructor(name, muscleGroup, sets){
        this.name = name;
        this.muscleGtoup = muscleGroup;
        this.sets = sets.map(s => new Set(s.repetitions, s.weight))
    }

    calculateEffort(){
        return this.sets.reduce((total, set) => total + set.calculateEffort(), 0);
    }
}

class Workout {
    constructor(userId, exercises){
        this.userId = userId;
        this.date = new Date();
        this.exercises = exercises.map(e => new ExerciseLog(e.name, e.muscleGroup, e.sets));
    }
}

module.exports = { Set, ExerciseLog, Workout };