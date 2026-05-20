package handlers

import (
	"context"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type eventRepo interface {
	List(ctx context.Context, topicID string, offset, limit int) ([]models.Event, int, error)
	Get(ctx context.Context, topicID, id string) (*models.Event, error)
}

type topicExistenceChecker interface {
	Exists(ctx context.Context, id string) (bool, error)
}

type EventHandler struct {
	repo       eventRepo
	topics     topicExistenceChecker
	dispatcher *events.Dispatcher
}

func NewEventHandler(repo eventRepo, topics topicExistenceChecker, d *events.Dispatcher) *EventHandler {
	return &EventHandler{repo: repo, topics: topics, dispatcher: d}
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topicId")
	exists, err := h.topics.Exists(r.Context(), topicID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if !exists {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Topic not found")
		return
	}
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), topicID, p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	topicID := r.PathValue("topicId")
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), topicID, id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Event not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}
