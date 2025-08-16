package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/tanq16/nojiko/internal/state"
)

func NewRouter(staticFS embed.FS, state *state.State) *http.ServeMux {
	mux := http.NewServeMux()
	apiHandler := NewAPIHandler(state)
	mux.HandleFunc("/api/bookmarks", apiHandler.GetBookmarks)
	mux.HandleFunc("/api/status-cards", apiHandler.GetStatusCards)
	mux.HandleFunc("/api/thumb-feeds", apiHandler.GetThumbFeeds)
	mux.HandleFunc("/api/header", apiHandler.GetHeader)
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticContent)))
	return mux
}
