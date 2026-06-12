package database

const (
	taskColumns = "id, title, description, completed, created_at, updated_at"

	queryGetAllTasks = `SELECT ` + taskColumns + ` FROM tasks ORDER BY created_at DESC`

	queryGetTaskByID = `SELECT ` + taskColumns + ` FROM tasks WHERE id = $1`

	queryCreateTask = `
INSERT INTO tasks (title, description, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + taskColumns

	queryUpdateTask = `
UPDATE tasks
SET title       = COALESCE($1, title),
	description = COALESCE($2, description),
	completed   = COALESCE($3, completed),
	updated_at  = $4
WHERE id = $5
RETURNING ` + taskColumns

	queryDeleteTask = `DELETE FROM tasks WHERE id = $1`
)
