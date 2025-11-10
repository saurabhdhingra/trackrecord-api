package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users         UserModel
	Exercises     ExerciseModel
	Workouts      WorkoutModel
	WorkoutItems  WorkoutItemModel
	WorkoutLogs   WorkoutLogModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:         UserModel{DB: db},
		Exercises:     ExerciseModel{DB: db},
		Workouts:      WorkoutModel{DB: db},
		WorkoutItems:  WorkoutItemModel{DB: db},
		WorkoutLogs:   WorkoutLogModel{DB: db},
	}
} 