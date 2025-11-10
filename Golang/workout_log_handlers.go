package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/trackrecord/api/internal/data"
	"github.com/trackrecord/api/internal/validator"
)

func (app *application) logWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorkoutID int64     `json:"workout_id"`
		Date      time.Time `json:"date"`
		Duration  int       `json:"duration"`
		Notes     string    `json:"notes"`
		Items     []struct {
			ExerciseID int64   `json:"exercise_id"`
			Sets       int     `json:"sets"`
			Reps       int     `json:"reps"`
			Weight     float64 `json:"weight"`
		} `json:"items"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	userIDStr := app.contextGetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.WorkoutID != 0, "workout_id", "must be provided")
	v.Check(input.Duration > 0, "duration", "must be greater than zero")
	v.Check(len(input.Items) > 0, "items", "must contain at least one item")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Verify workout belongs to user
	_, err = app.models.Workouts.Get(input.WorkoutID, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("workout_id", "workout not found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	workoutLog := &data.WorkoutLog{
		WorkoutID: input.WorkoutID,
		UserID:    userID,
		Date:      input.Date,
		Duration:  input.Duration,
		Notes:     input.Notes,
		Items:     make([]data.WorkoutLogItem, 0, len(input.Items)),
	}

	for i, item := range input.Items {
		v.Check(item.ExerciseID != 0, fmt.Sprintf("items.%d.exercise_id", i), "must be provided")
		v.Check(item.Sets > 0, fmt.Sprintf("items.%d.sets", i), "must be greater than zero")
		v.Check(item.Reps > 0, fmt.Sprintf("items.%d.reps", i), "must be greater than zero")
		v.Check(item.Weight >= 0, fmt.Sprintf("items.%d.weight", i), "must not be negative")

		workoutLog.Items = append(workoutLog.Items, data.WorkoutLogItem{
			ExerciseID: item.ExerciseID,
			Sets:       item.Sets,
			Reps:       item.Reps,
			Weight:     item.Weight,
		})
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.WorkoutLogs.Insert(workoutLog)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workout-logs/%d", workoutLog.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout_log": workoutLog}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getWorkoutLogHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	userIDStr := app.contextGetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	workoutLog, err := app.models.WorkoutLogs.Get(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout_log": workoutLog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listWorkoutLogsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := app.contextGetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "-date")

	workoutLogs, metadata, err := app.models.WorkoutLogs.GetAll(userID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout_logs": workoutLogs, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) generateProgressReportHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := app.contextGetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		ExerciseID int64 `json:"exercise_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		// Try query parameters if JSON body is not present
		qs := r.URL.Query()
		input.StartDate = app.readString(qs, "start_date", "")
		input.EndDate = app.readString(qs, "end_date", "")
		exerciseIDStr := app.readString(qs, "exercise_id", "")
		if exerciseIDStr != "" {
			exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
			if err == nil {
				input.ExerciseID = exerciseID
			}
		}
	}

	v := validator.New()

	var startDate, endDate time.Time
	if input.StartDate != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			v.AddError("start_date", "must be in format YYYY-MM-DD")
		}
	} else {
		// Default to 30 days ago
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if input.EndDate != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			v.AddError("end_date", "must be in format YYYY-MM-DD")
		}
	} else {
		// Default to today
		endDate = time.Now()
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// For now, we'll just return workout logs within the date range
	// In a real application, you'd analyze the data more deeply
	var filters data.Filters
	filters.Page = 1
	filters.PageSize = 100
	filters.Sort = "-date"

	workoutLogs, _, err := app.models.WorkoutLogs.GetAllBetweenDates(userID, startDate, endDate, input.ExerciseID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Calculate some basic stats
	report := map[string]interface{}{
		"start_date":      startDate.Format("2006-01-02"),
		"end_date":        endDate.Format("2006-01-02"),
		"workout_count":   len(workoutLogs),
		"total_duration":  0,
		"workouts":        workoutLogs,
		"exercise_stats":  make(map[string]interface{}),
	}

	exerciseStats := make(map[int64]map[string]interface{})
	totalDuration := 0

	for _, log := range workoutLogs {
		totalDuration += log.Duration

		// Collect exercise stats
		for _, item := range log.Items {
			if input.ExerciseID != 0 && item.ExerciseID != input.ExerciseID {
				continue
			}

			if _, exists := exerciseStats[item.ExerciseID]; !exists {
				exerciseStats[item.ExerciseID] = map[string]interface{}{
					"exercise_id":   item.ExerciseID,
					"exercise_name": item.Exercise.Name,
					"total_sets":    0,
					"total_reps":    0,
					"max_weight":    0.0,
					"workouts":      0,
				}
			}

			stats := exerciseStats[item.ExerciseID]
			stats["total_sets"] = stats["total_sets"].(int) + item.Sets
			stats["total_reps"] = stats["total_reps"].(int) + (item.Sets * item.Reps)
			stats["workouts"] = stats["workouts"].(int) + 1

			if item.Weight > stats["max_weight"].(float64) {
				stats["max_weight"] = item.Weight
			}
		}
	}

	report["total_duration"] = totalDuration
	report["exercise_stats"] = exerciseStats

	err = app.writeJSON(w, http.StatusOK, envelope{"report": report}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
} 