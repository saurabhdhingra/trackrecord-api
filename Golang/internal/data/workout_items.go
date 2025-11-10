package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type WorkoutItem struct {
	ID         int64     `json:"id"`
	WorkoutID  int64     `json:"workout_id"`
	ExerciseID int64     `json:"exercise_id"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	Weight     float64   `json:"weight"`
	Exercise   *Exercise `json:"exercise,omitempty"`
}

type WorkoutItemModel struct {
	DB *sql.DB
}

func (m WorkoutItemModel) Insert(item *WorkoutItem) error {
	query := `
		INSERT INTO workout_items (workout_id, exercise_id, sets, reps, weight)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		item.WorkoutID,
		item.ExerciseID,
		item.Sets,
		item.Reps,
		item.Weight,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
}

func (m WorkoutItemModel) Update(item *WorkoutItem) error {
	query := `
		UPDATE workout_items
		SET exercise_id = $1, sets = $2, reps = $3, weight = $4
		WHERE id = $5 AND workout_id = $6
		RETURNING id`

	args := []interface{}{
		item.ExerciseID,
		item.Sets,
		item.Reps,
		item.Weight,
		item.ID,
		item.WorkoutID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m WorkoutItemModel) Delete(id int64, workoutID int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM workout_items
		WHERE id = $1 AND workout_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, workoutID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m WorkoutItemModel) DeleteAllForWorkout(workoutID int64) error {
	query := `
		DELETE FROM workout_items
		WHERE workout_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, workoutID)
	return err
} 