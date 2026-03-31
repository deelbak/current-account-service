package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"current-account-service/internal/middleware"
	"current-account-service/internal/models"
)

// POST /applications
func (h *Handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	if middleware.RoleFromCtx(r) != "manager" {
		writeError(w, http.StatusForbidden, "only managers can create applications")
		return
	}

	var req models.CreateApplicationRequest
	if err := parseJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ClientID == 0 {
		writeError(w, http.StatusBadRequest, "client_id required")
		return
	}

	managerID := middleware.UserIDFromCtx(r)
	app, err := h.app.CreateApplication(r.Context(), req.ClientID, managerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("create failed: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, app)
}

// GET /applications/{id}
func (h *Handler) GetApplication(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	app, err := h.app.GetWithEvents(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "application not found")
		return
	}

	writeJSON(w, http.StatusOK, app)
}

// POST /applications/{id}/documents
func (h *Handler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	appID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Максимум 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file field required")
		return
	}
	defer file.Close()

	docType := r.FormValue("doc_type")
	if docType == "" {
		writeError(w, http.StatusBadRequest, "doc_type required")
		return
	}

	doc, err := h.app.UploadDocument(r.Context(), appID, docType, header.Filename, file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("upload failed: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, doc)
}

// POST /applications/{id}/submit
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	appID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	managerID := middleware.UserIDFromCtx(r)
	app, err := h.app.Submit(r.Context(), appID, managerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, app)
}

// POST /applications/{id}/approve
func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	appID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req models.SubmitRevisionRequest
	_ = parseJSON(r, &req) // comment опциональный

	headID := middleware.UserIDFromCtx(r)
	app, err := h.app.Approve(r.Context(), appID, headID, req.Comment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, app)
}

// POST /applications/{id}/revision
func (h *Handler) RequestRevision(w http.ResponseWriter, r *http.Request) {
	appID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req models.SubmitRevisionRequest
	if err := parseJSON(r, &req); err != nil || req.Comment == "" {
		writeError(w, http.StatusBadRequest, "comment required for revision")
		return
	}

	headID := middleware.UserIDFromCtx(r)
	app, err := h.app.RequestRevision(r.Context(), appID, headID, req.Comment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, app)
}

// POST /applications/{id}/reject  (бонус)
func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	appID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req models.SubmitRevisionRequest
	if err := parseJSON(r, &req); err != nil || req.Comment == "" {
		writeError(w, http.StatusBadRequest, "comment required for rejection")
		return
	}

	headID := middleware.UserIDFromCtx(r)
	app, err := h.app.Reject(r.Context(), appID, headID, req.Comment)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, app)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func parseID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(r.PathValue(param), 10, 64)
}
