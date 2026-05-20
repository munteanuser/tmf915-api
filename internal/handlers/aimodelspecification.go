package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type aiModelSpecRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.AiModelSpecification, int, error)
	Get(ctx context.Context, id string) (*models.AiModelSpecification, error)
	Create(ctx context.Context, in models.AiModelSpecificationCreate) (*models.AiModelSpecification, error)
	Update(ctx context.Context, id string, in models.AiModelSpecificationUpdate) (*models.AiModelSpecification, error)
	Delete(ctx context.Context, id string) error
}

type AiModelSpecificationHandler struct {
	repo       aiModelSpecRepo
	dispatcher *events.Dispatcher
}

func NewAiModelSpecificationHandler(repo aiModelSpecRepo, d *events.Dispatcher) *AiModelSpecificationHandler {
	return &AiModelSpecificationHandler{repo: repo, dispatcher: d}
}

func (h *AiModelSpecificationHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *AiModelSpecificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.AiModelSpecificationCreate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	if in.Name == "" {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", "name is required")
		return
	}
	item, err := h.repo.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AiModelSpecificationCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *AiModelSpecificationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiModelSpecification not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiModelSpecificationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.AiModelSpecificationUpdate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	item, err := h.repo.Update(r.Context(), id, in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiModelSpecification not found")
		return
	}
	h.dispatcher.Dispatch("AiModelSpecificationAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiModelSpecificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiModelSpecification not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AiModelSpecificationDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
