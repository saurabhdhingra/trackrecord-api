package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Health check endpoint
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Authentication endpoints
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.loginUserHandler)

	// Exercise endpoints
	router.HandlerFunc(http.MethodGet, "/v1/exercises", app.requireAuthenticatedUser(app.listExercisesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/exercises/:id", app.requireAuthenticatedUser(app.getExerciseHandler))

	// Workout endpoints
	router.HandlerFunc(http.MethodGet, "/v1/workouts", app.requireAuthenticatedUser(app.listWorkoutsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/workouts", app.requireAuthenticatedUser(app.createWorkoutHandler))
	router.HandlerFunc(http.MethodGet, "/v1/workouts/:id", app.requireAuthenticatedUser(app.getWorkoutHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/workouts/:id", app.requireAuthenticatedUser(app.updateWorkoutHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/workouts/:id", app.requireAuthenticatedUser(app.deleteWorkoutHandler))

	// Workout log endpoints
	router.HandlerFunc(http.MethodPost, "/v1/workout-logs", app.requireAuthenticatedUser(app.logWorkoutHandler))
	router.HandlerFunc(http.MethodGet, "/v1/workout-logs", app.requireAuthenticatedUser(app.listWorkoutLogsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/workout-logs/:id", app.requireAuthenticatedUser(app.getWorkoutLogHandler))

	// Report endpoints
	router.HandlerFunc(http.MethodGet, "/v1/reports/progress", app.requireAuthenticatedUser(app.generateProgressReportHandler))

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
} 