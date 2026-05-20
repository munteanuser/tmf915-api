package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type aiContractRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.AiContract, int, error)
	Get(ctx context.Context, id string) (*models.AiContract, error)
	Create(ctx context.Context, in models.AiContractCreate) (*models.AiContract, error)
	Update(ctx context.Context, id string, in models.AiContractUpdate) (*models.AiContract, error)
	Delete(ctx context.Context, id string) error
}

type AiContractHandler struct {
	repo       aiContractRepo
	dispatcher *events.Dispatcher
}

func NewAiContractHandler(repo aiContractRepo, d *events.Dispatcher) *AiContractHandler {
	return &AiContractHandler{repo: repo, dispatcher: d}
}

func (h *AiContractHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *AiContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.AiContractCreate
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
	h.dispatcher.Dispatch("AiContractCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *AiContractHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContract not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.AiContractUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContract not found")
		return
	}
	h.dispatcher.Dispatch("AiContractAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *AiContractHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "AiContract not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AiContractDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
