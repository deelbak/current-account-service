package handler

import (
	"net/http"

	"current-account-service/internal/middleware"
)

// GET /tasks/my
func (h *Handler) MyTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromCtx(r)

	tasks, err := h.task.ListMyTasks(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load tasks")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}
