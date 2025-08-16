package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/tanq16/nojiko/internal/api"
	"github.com/tanq16/nojiko/internal/config"
	"github.com/tanq16/nojiko/internal/state"
)

//go:embed static
var staticFS embed.FS

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appState := state.NewState(cfg)
	go appState.StartUpdateLoop()

	router := api.NewRouter(staticFS, appState)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
