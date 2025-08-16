package api

import (
	"encoding/json"
	"net/http"

	"github.com/tanq16/nojiko/internal/config"
	"github.com/tanq16/nojiko/internal/fetcher"
)

// APIHandler holds dependencies for API handlers.
type APIHandler struct {
	cfg *config.Config
}

// NewAPIHandler creates a new APIHandler.
func NewAPIHandler(cfg *config.Config) *APIHandler {
	return &APIHandler{cfg: cfg}
}

// GetBookmarks handles requests for bookmarks.
func (h *APIHandler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, h.cfg.Bookmarks)
}

// GetGitHubRepos handles requests for GitHub repository stats.
func (h *APIHandler) GetGitHubRepos(w http.ResponseWriter, r *http.Request) {
	// This currently returns mock data.
	// In a real implementation, it would fetch live data.
	repos := fetcher.GetGitHubStats(h.cfg.Github)
	respondWithJSON(w, http.StatusOK, repos)
}

// GetYouTubeVideos handles requests for the latest YouTube videos.
func (h *APIHandler) GetYouTubeVideos(w http.ResponseWriter, r *http.Request) {
	// This currently returns mock data.
	videos := fetcher.GetYouTubeVideos(h.cfg.Youtube)
	respondWithJSON(w, http.StatusOK, videos)
}

// GetHeader handles requests for header information like weather.
func (h *APIHandler) GetHeader(w http.ResponseWriter, r *http.Request) {
	// This currently returns mock data.
	headerInfo := fetcher.GetHeaderInfo()
	respondWithJSON(w, http.StatusOK, headerInfo)
}

// respondWithJSON is a helper to write JSON responses.
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
