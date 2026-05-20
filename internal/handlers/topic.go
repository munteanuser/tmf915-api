package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/models"
)

type topicRepo interface {
	List(ctx context.Context, offset, limit int) ([]models.Topic, int, error)
	Get(ctx context.Context, id string) (*models.Topic, error)
	Create(ctx context.Context, in models.TopicCreate) (*models.Topic, error)
	Update(ctx context.Context, id string, in models.TopicUpdate) (*models.Topic, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

type TopicHandler struct {
	repo       topicRepo
	dispatcher *events.Dispatcher
}

func NewTopicHandler(repo topicRepo, d *events.Dispatcher) *TopicHandler {
	return &TopicHandler{repo: repo, dispatcher: d}
}

func (h *TopicHandler) List(w http.ResponseWriter, r *http.Request) {
	p := ParseListParams(r)
	items, total, err := h.repo.List(r.Context(), p.Offset, p.Limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	SetPaginationHeaders(w, len(items), total)
	WriteJSON(w, http.StatusOK, items)
}

func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.TopicCreate
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
	h.dispatcher.Dispatch("TopicCreateEvent", item)
	WriteJSON(w, http.StatusCreated, item)
}

func (h *TopicHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Topic not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *TopicHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in models.TopicUpdate
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
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Topic not found")
		return
	}
	h.dispatcher.Dispatch("TopicChangeEvent", item)
	WriteJSON(w, http.StatusOK, item)
}

func (h *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Topic not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	h.dispatcher.Dispatch("TopicDeleteEvent", map[string]string{"id": id})
	w.WriteHeader(http.StatusNoContent)
}
