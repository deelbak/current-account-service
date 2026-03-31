package handler

import (
	"net/http"
)

// GET /clients/{iin}
func (h *Handler) GetClient(w http.ResponseWriter, r *http.Request) {
	iin := r.PathValue("iin")
	if len(iin) != 12 {
		writeError(w, http.StatusBadRequest, "IIN must be 12 characters")
		return
	}

	client, err := h.client.GetByIIN(r.Context(), iin)
	if err != nil {
		writeError(w, http.StatusNotFound, "client not found")
		return
	}

	writeJSON(w, http.StatusOK, client)
}
