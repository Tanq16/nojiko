package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/tanq16/nojiko/internal/config"
)

// NewRouter sets up the application's HTTP routes.
func NewRouter(staticFS embed.FS, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	apiHandler := NewAPIHandler(cfg)

	// Register API endpoints.
	mux.HandleFunc("/api/bookmarks", apiHandler.GetBookmarks)
	mux.HandleFunc("/api/github", apiHandler.GetGitHubRepos)
	mux.HandleFunc("/api/youtube", apiHandler.GetYouTubeVideos)
	mux.HandleFunc("/api/header", apiHandler.GetHeader)

	// Create a file system from the embedded static assets.
	// The "static" directory is stripped from the path.
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		// This should not happen with a valid embed.
		panic(err)
	}

	// Serve the static files (HTML, JS, etc.).
	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	return mux
}
