package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type monitorRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.Monitor, int, error)
	Get(ctx context.Context, id string) (*models.Monitor, error)
	Create(ctx context.Context, in models.MonitorCreate) (*models.Monitor, error)
	Update(ctx context.Context, id string, in models.MonitorUpdate) (*models.Monitor, error)
	Delete(ctx context.Context, id string) error
}

type MonitorHandler struct {
	repo       monitorRepo
	dispatcher *events.Dispatcher
}

func NewMonitorHandler(repo monitorRepo, d *events.Dispatcher) *MonitorHandler {
	return &MonitorHandler{repo: repo, dispatcher: d}
}

func (h *MonitorHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *MonitorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.MonitorCreate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	item, err := h.repo.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("MonitorCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *MonitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Monitor not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *MonitorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.MonitorUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Monitor not found")
		return
	}
	h.dispatcher.Dispatch("MonitorAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *MonitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Monitor not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("MonitorDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
