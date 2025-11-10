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

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Schedule    time.Time `json:"schedule"`
		Items       []struct {
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

	workout := &data.Workout{
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
		Schedule:    input.Schedule,
		Items:       make([]data.WorkoutItem, 0, len(input.Items)),
	}

	for _, item := range input.Items {
		workout.Items = append(workout.Items, data.WorkoutItem{
			ExerciseID: item.ExerciseID,
			Sets:       item.Sets,
			Reps:       item.Reps,
			Weight:     item.Weight,
		})
	}

	v := validator.New()

	if validateWorkout(v, workout); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Workouts.Insert(workout)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Insert workout items
	for i := range workout.Items {
		workout.Items[i].WorkoutID = workout.ID
		err = app.models.WorkoutItems.Insert(&workout.Items[i])
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workouts/%d", workout.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getWorkoutHandler(w http.ResponseWriter, r *http.Request) {
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

	workout, err := app.models.Workouts.Get(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
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

	workout, err := app.models.Workouts.Get(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		Schedule    *time.Time `json:"schedule"`
		Items       []struct {
			ID         *int64   `json:"id"`
			ExerciseID int64    `json:"exercise_id"`
			Sets       int      `json:"sets"`
			Reps       int      `json:"reps"`
			Weight     float64  `json:"weight"`
		} `json:"items"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		workout.Name = *input.Name
	}

	if input.Description != nil {
		workout.Description = *input.Description
	}

	if input.Schedule != nil {
		workout.Schedule = *input.Schedule
	}

	v := validator.New()

	if validateWorkout(v, workout); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Workouts.Update(workout)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Items != nil {
		// Delete all existing workout items
		err = app.models.WorkoutItems.DeleteAllForWorkout(workout.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Insert new workout items
		workout.Items = make([]data.WorkoutItem, 0, len(input.Items))
		for _, item := range input.Items {
			workoutItem := data.WorkoutItem{
				WorkoutID:  workout.ID,
				ExerciseID: item.ExerciseID,
				Sets:       item.Sets,
				Reps:       item.Reps,
				Weight:     item.Weight,
			}

			if item.ID != nil {
				workoutItem.ID = *item.ID
			}

			err = app.models.WorkoutItems.Insert(&workoutItem)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			workout.Items = append(workout.Items, workoutItem)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = app.models.Workouts.Delete(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "workout successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
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
	input.Filters.Sort = app.readString(qs, "sort", "created_at")

	workouts, metadata, err := app.models.Workouts.GetAll(userID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workouts": workouts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func validateWorkout(v *validator.Validator, workout *data.Workout) bool {
	v.Check(workout.Name != "", "name", "must be provided")
	v.Check(len(workout.Name) <= 100, "name", "must not be more than 100 bytes long")

	v.Check(workout.Description != "", "description", "must be provided")
	v.Check(len(workout.Description) <= 1000, "description", "must not be more than 1000 bytes long")

	if !workout.Schedule.IsZero() {
		v.Check(workout.Schedule.After(time.Now()), "schedule", "must be in the future")
	}

	if len(workout.Items) > 0 {
		for i, item := range workout.Items {
			v.Check(item.ExerciseID != 0, fmt.Sprintf("items.%d.exercise_id", i), "must be provided")
			v.Check(item.Sets > 0, fmt.Sprintf("items.%d.sets", i), "must be greater than zero")
			v.Check(item.Reps > 0, fmt.Sprintf("items.%d.reps", i), "must be greater than zero")
			v.Check(item.Weight >= 0, fmt.Sprintf("items.%d.weight", i), "must not be negative")
		}
	}

	return v.Valid()
} 