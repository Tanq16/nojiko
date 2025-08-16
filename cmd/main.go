package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/tanq16/nojiko/internal/api"
	"github.com/tanq16/nojiko/internal/config"
)

//go:embed static
var staticFS embed.FS

// main is the application entry point.
func main() {
	// Load application configuration from the YAML file.
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create a new router and register handlers.
	router := api.NewRouter(staticFS, cfg)

	// Start the HTTP server.
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
