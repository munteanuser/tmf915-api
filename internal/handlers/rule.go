package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type ruleRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.Rule, int, error)
	Get(ctx context.Context, id string) (*models.Rule, error)
	Create(ctx context.Context, in models.RuleCreate) (*models.Rule, error)
	Update(ctx context.Context, id string, in models.RuleUpdate) (*models.Rule, error)
	Delete(ctx context.Context, id string) error
}

type RuleHandler struct {
	repo       ruleRepo
	dispatcher *events.Dispatcher
}

func NewRuleHandler(repo ruleRepo, d *events.Dispatcher) *RuleHandler {
	return &RuleHandler{repo: repo, dispatcher: d}
}

func (h *RuleHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *RuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.RuleCreate
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
	h.dispatcher.Dispatch("RuleCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *RuleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Rule not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *RuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.RuleUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Rule not found")
		return
	}
	h.dispatcher.Dispatch("RuleAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *RuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Rule not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("RuleDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
