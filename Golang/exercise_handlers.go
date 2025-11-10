package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/trackrecord/api/internal/data"
	"github.com/trackrecord/api/internal/validator"
)

func (app *application) listExercisesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Category    string
		MuscleGroup string
		data.Filters
	}

	v := validator.New()
	
	qs := r.URL.Query()
	
	input.Category = app.readString(qs, "category", "")
	input.MuscleGroup = app.readString(qs, "muscle_group", "")
	
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "name")
	
	exercises, metadata, err := app.models.Exercises.GetAll(input.Category, input.MuscleGroup, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"exercises": exercises, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
	exercise, err := app.models.Exercises.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
} 