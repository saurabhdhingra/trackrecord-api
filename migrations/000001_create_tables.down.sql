-- Drop tables in reverse order of creation (due to foreign key constraints)
DROP TABLE IF EXISTS workout_log_items;
DROP TABLE IF EXISTS workout_logs;
DROP TABLE IF EXISTS workout_items;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS users; 