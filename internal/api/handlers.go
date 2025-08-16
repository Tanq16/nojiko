package api

import (
	"encoding/json"
	"net/http"

	"github.com/tanq16/nojiko/internal/state"
)

type APIHandler struct {
	state *state.State
}

func NewAPIHandler(state *state.State) *APIHandler {
	return &APIHandler{state: state}
}

func (h *APIHandler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.state.GetBookmarks())
}

func (h *APIHandler) GetStatusCards(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.state.GetStatusCards())
}

func (h *APIHandler) GetThumbFeeds(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.state.GetThumbFeeds())
}

func (h *APIHandler) GetHeader(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.state.GetHeader())
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
