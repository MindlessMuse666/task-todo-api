// TODO(MindlessMuse): задокументировать модуль и методы
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MindlessMuse666/task-todo-api/internal/models"
	"github.com/jmoiron/sqlx"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (s *TaskStore) GetAll() ([]models.Task, error) {
	tasks := []models.Task{}

	// TODO(MindlessMuse): вынести в отдельный скрипт
	query := `
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
ORDER BY created_at DESC;
`

	if err := s.db.Select(&tasks, query); err != nil {
		return nil, fmt.Errorf("Ошибка получения всех задач: %w", err)
	}

	return tasks, nil
}

func (s *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	// TODO(MindlessMuse): вынести в отдельный скрипт
	query := `
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
WHERE id = $1;
`

	if err := s.db.Get(&task, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача с id=%d не найдена", id)
		}

		return nil, fmt.Errorf("Ошибка получения задачи с id=%d: %w", id, err)
	}

	return &task, nil
}

// Create создает новую задачу в базе данных и возвращает ее.
func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task

	// TODO(MindlessMuse): вынести в отдельный скрипт
	query := `
INSERT INTO tasks(title, description, completed, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
returning id, title, description, completed, created_at, updated_at;`

	now := time.Now()

	if err := s.db.QueryRowx(query, input.Title, input.Description, input.Completed, now, now).StructScan(&task); err != nil {
		return nil, fmt.Errorf("Ошибка создания задачи: %w", err)
	}

	return &task, nil
}
