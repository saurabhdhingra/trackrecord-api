package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Workout struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	UserID      int64        `json:"user_id"`
	Schedule    time.Time    `json:"schedule,omitempty"`
	Items       []WorkoutItem `json:"items,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type WorkoutModel struct {
	DB *sql.DB
}

func (m WorkoutModel) Insert(workout *Workout) error {
	query := `
		INSERT INTO workouts (name, description, user_id, schedule, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		workout.Name,
		workout.Description,
		workout.UserID,
		workout.Schedule,
		time.Now(),
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&workout.ID)
}

func (m WorkoutModel) Get(id int64, userID int64) (*Workout, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT w.id, w.name, w.description, w.user_id, w.schedule, w.created_at, w.updated_at
		FROM workouts w
		WHERE w.id = $1 AND w.user_id = $2`

	var workout Workout

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&workout.ID,
		&workout.Name,
		&workout.Description,
		&workout.UserID,
		&workout.Schedule,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Get workout items
	workout.Items, err = GetWorkoutItems(m.DB, workout.ID)
	if err != nil {
		return nil, err
	}

	return &workout, nil
}

func (m WorkoutModel) Update(workout *Workout) error {
	query := `
		UPDATE workouts
		SET name = $1, description = $2, schedule = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6
		RETURNING id`

	args := []interface{}{
		workout.Name,
		workout.Description,
		workout.Schedule,
		time.Now(),
		workout.ID,
		workout.UserID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&workout.ID)
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

func (m WorkoutModel) Delete(id int64, userID int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// Start a transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete workout items first
	_, err = tx.ExecContext(ctx, "DELETE FROM workout_items WHERE workout_id = $1", id)
	if err != nil {
		return err
	}

	// Delete workout
	result, err := tx.ExecContext(ctx, "DELETE FROM workouts WHERE id = $1 AND user_id = $2", id, userID)
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

	// Commit the transaction
	return tx.Commit()
}

func (m WorkoutModel) GetAll(userID int64, filters Filters) ([]*Workout, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, user_id, schedule, created_at, updated_at
		FROM workouts
		WHERE user_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{userID, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	workouts := []*Workout{}

	for rows.Next() {
		var workout Workout

		err := rows.Scan(
			&totalRecords,
			&workout.ID,
			&workout.Name,
			&workout.Description,
			&workout.UserID,
			&workout.Schedule,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Get workout items
		workout.Items, err = GetWorkoutItems(m.DB, workout.ID)
		if err != nil {
			return nil, Metadata{}, err
		}

		workouts = append(workouts, &workout)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return workouts, metadata, nil
}

func GetWorkoutItems(db *sql.DB, workoutID int64) ([]WorkoutItem, error) {
	query := `
		SELECT wi.id, wi.workout_id, wi.exercise_id, wi.sets, wi.reps, wi.weight, 
			   e.name, e.description, e.category, e.muscle_group
		FROM workout_items wi
		JOIN exercises e ON wi.exercise_id = e.id
		WHERE wi.workout_id = $1
		ORDER BY wi.id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []WorkoutItem{}

	for rows.Next() {
		var item WorkoutItem
		var exercise Exercise

		err := rows.Scan(
			&item.ID,
			&item.WorkoutID,
			&item.ExerciseID,
			&item.Sets,
			&item.Reps,
			&item.Weight,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.MuscleGroup,
		)
		if err != nil {
			return nil, err
		}

		item.Exercise = &exercise
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
} 