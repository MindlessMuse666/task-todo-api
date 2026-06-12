// Package database реализует слой работы с PostgreSQL-хранилищем задач.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MindlessMuse666/task-todo-api/internal/models"
	"github.com/jmoiron/sqlx"
)

// ErrNotFound возвращается, когда искомая запись отсутствует.
var ErrNotFound = errors.New("record not found")

// TaskStore хранилище задач.
type TaskStore struct {
	db *sqlx.DB
}

// NewTaskStore создаёт новый экземпляр TaskStore.
func NewTaskStore(db *sqlx.DB) *TaskStore {
	return &TaskStore{db: db}
}

// GetAll возвращает все задачи, упорядоченные по дате создания (сначала новые).
func (s *TaskStore) GetAll() ([]models.Task, error) {
	tasks := []models.Task{}

	if err := s.db.Select(&tasks, queryGetAllTasks); err != nil {
		return nil, fmt.Errorf("Ошибка получения всех задач: %w", err)
	}

	return tasks, nil
}

// GetByID возвращает задачу по идентификатору.
// Возвращает ErrNotFound, если задача не найдена.
func (s *TaskStore) GetByID(id int) (*models.Task, error) {
	var task models.Task

	if err := s.db.Get(&task, queryGetTaskByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача с id=%d не найдена: %w", id, ErrNotFound)
		}

		return nil, fmt.Errorf("Ошибка получения задачи с id=%d: %w", id, err)
	}

	return &task, nil
}

// Create создает новую задачу и возвращает ее.
func (s *TaskStore) Create(input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task
	now := time.Now()

	if err := s.db.QueryRowx(queryCreateTask,
		input.Title,
		input.Description,
		input.Completed,
		now,
		now,
	).StructScan(&task); err != nil {
		return nil, fmt.Errorf("Ошибка создания задачи: %w", err)
	}

	return &task, nil
}

// Update частично обновляет задачу по идентификатору.
// Если поля в input равны nil, соответствующие колонки остаются без изменений.
// Возвращает ErrNotFound, если задача с указанным id отсутствует.
func (s *TaskStore) Update(id int, input models.UpdateTaskInput) (*models.Task, error) {
	var task models.Task
	now := time.Now()

	err := s.db.QueryRowx(queryUpdateTask,
		input.Title,
		input.Description,
		input.Completed,
		now,
		id,
	).StructScan(&task)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Задача с id=%d не найдена: %w", id, ErrNotFound)
		}

		return nil, fmt.Errorf("Ошибка обновления задачи с id=%d: %w", id, err)
	}

	return &task, nil
}

// Delete удаляет задачу по идентификатору.
// Возвращает ErrNotFound, если задача не найдена.
func (s *TaskStore) Delete(id int) error {
	result, err := s.db.Exec(queryDeleteTask, id)
	if err != nil {
		return fmt.Errorf("Ошибка удаления задачи с id=%d: %w", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Ошибка получения количество удаленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("Задача с id=%d не найдена: %w", id, ErrNotFound)
	}

	return nil
}
