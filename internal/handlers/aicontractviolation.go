package handlers

import (
	"context"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type aiContractViolationRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.AiContractViolation, int, error)
	Get(ctx context.Context, id string) (*models.AiContractViolation, error)
	Create(ctx context.Context, in models.AiContractViolationCreate) (*models.AiContractViolation, error)
}

type AiContractViolationHandler struct {
	repo       aiContractViolationRepo
	dispatcher *events.Dispatcher
}

func NewAiContractViolationHandler(repo aiContractViolationRepo, d *events.Dispatcher) *AiContractViolationHandler {
	return &AiContractViolationHandler{repo: repo, dispatcher: d}
}

func (h *AiContractViolationHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *AiContractViolationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.AiContractViolationCreate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	item, err := h.repo.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AiContractViolationCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *AiContractViolationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContractViolation not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}
