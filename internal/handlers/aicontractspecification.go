package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type aiContractSpecRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.AiContractSpecification, int, error)
	Get(ctx context.Context, id string) (*models.AiContractSpecification, error)
	Create(ctx context.Context, in models.AiContractSpecificationCreate) (*models.AiContractSpecification, error)
	Update(ctx context.Context, id string, in models.AiContractSpecificationUpdate) (*models.AiContractSpecification, error)
	Delete(ctx context.Context, id string) error
}

type AiContractSpecificationHandler struct {
	repo       aiContractSpecRepo
	dispatcher *events.Dispatcher
}

func NewAiContractSpecificationHandler(repo aiContractSpecRepo, d *events.Dispatcher) *AiContractSpecificationHandler {
	return &AiContractSpecificationHandler{repo: repo, dispatcher: d}
}

func (h *AiContractSpecificationHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *AiContractSpecificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.AiContractSpecificationCreate
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
	h.dispatcher.Dispatch("AiContractSpecificationCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *AiContractSpecificationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContractSpecification not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiContractSpecificationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.AiContractSpecificationUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContractSpecification not found")
		return
	}
	h.dispatcher.Dispatch("AiContractSpecificationAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiContractSpecificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContractSpecification not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AiContractSpecificationDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
