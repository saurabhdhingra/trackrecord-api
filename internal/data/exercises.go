package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Exercise struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	MuscleGroup string    `json:"muscle_group"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExerciseModel struct {
	DB *sql.DB
}

func (m ExerciseModel) Insert(exercise *Exercise) error {
	query := `
		INSERT INTO exercises (name, description, category, muscle_group, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		exercise.Name,
		exercise.Description,
		exercise.Category,
		exercise.MuscleGroup,
		time.Now(),
		time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID)
}

func (m ExerciseModel) Get(id int64) (*Exercise, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, description, category, muscle_group, created_at, updated_at
		FROM exercises
		WHERE id = $1`

	var exercise Exercise

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Category,
		&exercise.MuscleGroup,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &exercise, nil
}

func (m ExerciseModel) Update(exercise *Exercise) error {
	query := `
		UPDATE exercises
		SET name = $1, description = $2, category = $3, muscle_group = $4, updated_at = $5
		WHERE id = $6
		RETURNING id`

	args := []interface{}{
		exercise.Name,
		exercise.Description,
		exercise.Category,
		exercise.MuscleGroup,
		time.Now(),
		exercise.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID)
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

func (m ExerciseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM exercises
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m ExerciseModel) GetAll(category string, muscleGroup string, filters Filters) ([]*Exercise, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, category, muscle_group, created_at, updated_at
		FROM exercises
		WHERE (LOWER(category) = LOWER($1) OR $1 = '')
		AND (LOWER(muscle_group) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{category, muscleGroup, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	exercises := []*Exercise{}

	for rows.Next() {
		var exercise Exercise

		err := rows.Scan(
			&totalRecords,
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Category,
			&exercise.MuscleGroup,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		exercises = append(exercises, &exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return exercises, metadata, nil
} 