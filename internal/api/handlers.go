package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/tanq16/nojiko/internal/state"
)

type APIHandler struct {
	state *state.State
}

func NewAPIHandler(state *state.State) *APIHandler {
	return &APIHandler{state: state}
}

func (h *APIHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	configBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		http.Error(w, "Failed to read config file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.Write(configBytes)
}

func (h *APIHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var temp any
	if err := yaml.Unmarshal(body, &temp); err != nil {
		http.Error(w, "Invalid YAML format: "+err.Error(), http.StatusBadRequest)
		return
	}
	// backup the old config before writing
	if err := os.Rename("config.yaml", "config.yaml.bak"); err != nil {
		// still proceed on error
	}
	err = os.WriteFile("config.yaml", body, 0644)
	if err != nil {
		http.Error(w, "Failed to write config file", http.StatusInternalServerError)
		// restore backup if write fails
		os.Rename("config.yaml.bak", "config.yaml")
		return
	}
	err = h.state.ReloadConfig()
	if err != nil {
		http.Error(w, "Failed to reload new config: "+err.Error(), http.StatusInternalServerError)
		// restore backup if invalid config
		os.Rename("config.yaml.bak", "config.yaml")
		h.state.ReloadConfig()
		return
	}
	w.WriteHeader(http.StatusOK)
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

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
