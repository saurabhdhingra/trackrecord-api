class User {
    constructor(id, email, birthdate, measurements = {}, exercises = {}){
        this.id = id;
        this.email = email;
        this.birthdate = birthdate;
        this.measurements = measurements;
        this.exercises = exercises;
    }
}

class Measurement {
    constructor(weight, height, leftArmGirth, rightArmGirth, leftThighGirth, rightThighGirth, waist, shoulders, chest){
        this.weight = weight;
        this.height = height;
        this.leftArmGirth = leftArmGirth;
        this.rightArmGirth = rightArmGirth;
        this.leftThighGirth = leftThighGirth;
        this.rightThighGirth = rightThighGirth;
        this.waist = waist;
        this.shoulders = shoulders;
        this.chest = chest;
        this.date = new Date();
    }
}

MediaSourceHandle.exports = {User, Measurement};