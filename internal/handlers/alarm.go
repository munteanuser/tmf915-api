package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type alarmRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.Alarm, int, error)
	Get(ctx context.Context, id string) (*models.Alarm, error)
	Create(ctx context.Context, in models.AlarmCreate) (*models.Alarm, error)
	Update(ctx context.Context, id string, in models.AlarmUpdate) (*models.Alarm, error)
	Delete(ctx context.Context, id string) error
}

type AlarmHandler struct {
	repo       alarmRepo
	dispatcher *events.Dispatcher
}

func NewAlarmHandler(repo alarmRepo, d *events.Dispatcher) *AlarmHandler {
	return &AlarmHandler{repo: repo, dispatcher: d}
}

func (h *AlarmHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *AlarmHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.AlarmCreate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	item, err := h.repo.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AlarmCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *AlarmHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Alarm not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *AlarmHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.AlarmUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Alarm not found")
		return
	}
	h.dispatcher.Dispatch("AlarmAttributeValueChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *AlarmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Alarm not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("AlarmDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
