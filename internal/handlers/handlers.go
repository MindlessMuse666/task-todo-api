// TODO(MindlessMuse): задокументировать модуль и докстринги
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MindlessMuse666/task-todo-api/internal/database"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}
