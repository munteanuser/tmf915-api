package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tmf915-api/internal/models"
)

type hubRepo interface {
	List(ctx context.Context) ([]models.Hub, error)
	Get(ctx context.Context, id string) (*models.Hub, error)
	Create(ctx context.Context, in models.HubCreate) (*models.Hub, error)
	Delete(ctx context.Context, id string) error
}

type listenerRepo interface {
	Create(ctx context.Context, in models.EventSubscriptionInput) (*models.EventSubscription, error)
	Delete(ctx context.Context, id string) error
}

type HubHandler struct {
	repo     hubRepo
	listener listenerRepo
}

func NewHubHandler(repo hubRepo, listener listenerRepo) *HubHandler {
	return &HubHandler{repo: repo, listener: listener}
}

func (h *HubHandler) ListHubs(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, items)
}

func (h *HubHandler) CreateHub(w http.ResponseWriter, r *http.Request) {
	var in models.HubCreate
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	if in.Callback == "" {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", "callback is required")
		return
	}
	item, err := h.repo.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, item)
}

func (h *HubHandler) GetHub(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.repo.Get(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	if item == nil {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Hub not found")
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *HubHandler) DeleteHub(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.repo.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Hub not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HubHandler) RegisterListener(w http.ResponseWriter, r *http.Request) {
	var in models.EventSubscriptionInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", err.Error())
		return
	}
	if in.Callback == "" {
		WriteError(w, http.StatusBadRequest, "400", "Bad Request", "callback is required")
		return
	}
	item, err := h.listener.Create(r.Context(), in)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, item)
}

func (h *HubHandler) UnregisterListener(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.listener.Delete(r.Context(), id)
	if err == sql.ErrNoRows {
		WriteError(w, http.StatusNotFound, "404", "Not Found", "Listener not found")
		return
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "500", "Internal Server Error", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
