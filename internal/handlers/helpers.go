package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tmf915-api/internal/models"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, reason, message string) {
	msg := message
	WriteJSON(w, status, models.Error{
		Code:    code,
		Reason:  reason,
		Message: &msg,
	})
}

func DecodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

type ListParams struct {
	Fields string
	Offset int
	Limit  int
}

func ParseListParams(r *http.Request) ListParams {
	q := r.URL.Query()
	p := ListParams{
		Fields: q.Get("fields"),
		Limit:  100,
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			p.Offset = n
		}
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 1000 {
			p.Limit = n
		}
	}
	return p
}

func SetPaginationHeaders(w http.ResponseWriter, resultCount, totalCount int) {
	w.Header().Set("X-Result-Count", strconv.Itoa(resultCount))
	w.Header().Set("X-Total-Count", strconv.Itoa(totalCount))
}
