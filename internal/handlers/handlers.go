// Package handlers содержит HTTP-обработчики для REST API задач.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/MindlessMuse666/task-todo-api/internal/database"
	"github.com/MindlessMuse666/task-todo-api/internal/models"
)

// Handlers хранит зависимости обработчиков (слой доступа к данным).
type Handlers struct {
	store *database.TaskStore
}

// NewHandlers создаёт новый экземпляр Handlers.
func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

// GetAllTasks возвращает список всех задач.
func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

// GetTask возвращает задачу по идентификатору из URL.
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractTaskID(r.URL.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// CreateTask создаёт новую задачу из JSON-тела запроса.
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные отправленные данные")
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок задачи обязателен")
		return
	}

	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка создания задачи")
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

// UpdateTask обновляет существующую задачу по ID.
func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractTaskID(r.URL.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные отправленные данные")
		return
	}
	defer r.Body.Close()

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок задачи не может быть пустым")
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Ошибка обновления задачи")
		}

		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// DeleteTask удаляет задачу по ID.
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractTaskID(r.URL.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	if err = h.store.Delete(id); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, "Задача не найдена")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Ошибка удаления задачи")
		}

		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

/* ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ */

// respondWithJSON устанавливает заголовки и сериализует payload в JSON.
func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError отправляет JSON-ответ с ошибкой.
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// extractTaskID извлекает ID задачи из URL.
// Ожидается формат "/tasks/{id}".
func extractTaskID(path string) (int, error) {
	const prefix = "/tasks/"
	idStr := strings.TrimPrefix(path, prefix)
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, errors.New("invalid id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
