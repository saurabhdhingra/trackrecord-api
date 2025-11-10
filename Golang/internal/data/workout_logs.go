package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type WorkoutLog struct {
	ID          int64     `json:"id"`
	WorkoutID   int64     `json:"workout_id"`
	UserID      int64     `json:"user_id"`
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"` // in minutes
	Notes       string    `json:"notes"`
	WorkoutName string    `json:"workout_name,omitempty"`
	Items       []WorkoutLogItem `json:"items,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type WorkoutLogItem struct {
	ID         int64     `json:"id"`
	LogID      int64     `json:"log_id"`
	ExerciseID int64     `json:"exercise_id"`
	Sets       int       `json:"sets"`
	Reps       int       `json:"reps"`
	Weight     float64   `json:"weight"`
	Exercise   *Exercise `json:"exercise,omitempty"`
}

type WorkoutLogModel struct {
	DB *sql.DB
}

func (m WorkoutLogModel) Insert(log *WorkoutLog) error {
	// Start a transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert workout log
	query := `
		INSERT INTO workout_logs (workout_id, user_id, date, duration, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	args := []interface{}{
		log.WorkoutID,
		log.UserID,
		log.Date,
		log.Duration,
		log.Notes,
		time.Now(),
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(&log.ID)
	if err != nil {
		return err
	}

	// Insert workout log items
	for i := range log.Items {
		query = `
			INSERT INTO workout_log_items (log_id, exercise_id, sets, reps, weight)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`

		args = []interface{}{
			log.ID,
			log.Items[i].ExerciseID,
			log.Items[i].Sets,
			log.Items[i].Reps,
			log.Items[i].Weight,
		}

		err = tx.QueryRowContext(ctx, query, args...).Scan(&log.Items[i].ID)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (m WorkoutLogModel) Get(id int64, userID int64) (*WorkoutLog, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT wl.id, wl.workout_id, wl.user_id, wl.date, wl.duration, wl.notes, wl.created_at, w.name
		FROM workout_logs wl
		JOIN workouts w ON wl.workout_id = w.id
		WHERE wl.id = $1 AND wl.user_id = $2`

	var log WorkoutLog

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&log.ID,
		&log.WorkoutID,
		&log.UserID,
		&log.Date,
		&log.Duration,
		&log.Notes,
		&log.CreatedAt,
		&log.WorkoutName,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Get workout log items
	log.Items, err = GetWorkoutLogItems(m.DB, log.ID)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (m WorkoutLogModel) GetAll(userID int64, filters Filters) ([]*WorkoutLog, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), wl.id, wl.workout_id, wl.user_id, wl.date, wl.duration, wl.notes, wl.created_at, w.name
		FROM workout_logs wl
		JOIN workouts w ON wl.workout_id = w.id
		WHERE wl.user_id = $1
		ORDER BY %s %s, wl.id ASC
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
	logs := []*WorkoutLog{}

	for rows.Next() {
		var log WorkoutLog

		err := rows.Scan(
			&totalRecords,
			&log.ID,
			&log.WorkoutID,
			&log.UserID,
			&log.Date,
			&log.Duration,
			&log.Notes,
			&log.CreatedAt,
			&log.WorkoutName,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Get workout log items
		log.Items, err = GetWorkoutLogItems(m.DB, log.ID)
		if err != nil {
			return nil, Metadata{}, err
		}

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return logs, metadata, nil
}

func (m WorkoutLogModel) GetAllBetweenDates(userID int64, startDate, endDate time.Time, exerciseID int64, filters Filters) ([]*WorkoutLog, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), wl.id, wl.workout_id, wl.user_id, wl.date, wl.duration, wl.notes, wl.created_at, w.name
		FROM workout_logs wl
		JOIN workouts w ON wl.workout_id = w.id
		WHERE wl.user_id = $1
		AND wl.date >= $2
		AND wl.date <= $3
		%s
		ORDER BY %s %s, wl.id ASC
		LIMIT $5 OFFSET $6`, exerciseIDFilter(exerciseID), filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{userID, startDate, endDate}

	if exerciseID > 0 {
		// Add exerciseID to args
		args = append(args, exerciseID)
	}

	args = append(args, filters.limit(), filters.offset())

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	logs := []*WorkoutLog{}

	for rows.Next() {
		var log WorkoutLog

		err := rows.Scan(
			&totalRecords,
			&log.ID,
			&log.WorkoutID,
			&log.UserID,
			&log.Date,
			&log.Duration,
			&log.Notes,
			&log.CreatedAt,
			&log.WorkoutName,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Get workout log items
		log.Items, err = GetWorkoutLogItems(m.DB, log.ID)
		if err != nil {
			return nil, Metadata{}, err
		}

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return logs, metadata, nil
}

func exerciseIDFilter(exerciseID int64) string {
	if exerciseID > 0 {
		return `AND wl.id IN (
			SELECT log_id 
			FROM workout_log_items 
			WHERE exercise_id = $4
		)`
	}
	return ""
}

func GetWorkoutLogItems(db *sql.DB, logID int64) ([]WorkoutLogItem, error) {
	query := `
		SELECT wli.id, wli.log_id, wli.exercise_id, wli.sets, wli.reps, wli.weight, 
			   e.name, e.description, e.category, e.muscle_group
		FROM workout_log_items wli
		JOIN exercises e ON wli.exercise_id = e.id
		WHERE wli.log_id = $1
		ORDER BY wli.id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []WorkoutLogItem{}

	for rows.Next() {
		var item WorkoutLogItem
		var exercise Exercise

		err := rows.Scan(
			&item.ID,
			&item.LogID,
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