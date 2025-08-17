package state

import (
	"log"
	"sync"
	"time"

	"github.com/tanq16/nojiko/internal/config"
	"github.com/tanq16/nojiko/internal/fetcher"
	"github.com/tanq16/nojiko/internal/models"
)

type State struct {
	mu          sync.RWMutex
	cfg         *config.Config
	header      *models.HeaderInfo
	statusCards []models.StatusCardSection
	thumbFeeds  []models.ThumbFeedSection
	bookmarks   []config.BookmarkCategory
}

func NewState(cfg *config.Config) *State {
	s := &State{
		cfg:       cfg,
		bookmarks: cfg.Bookmarks,
	}
	s.updateState()
	return s
}

func (s *State) updateState() {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println("Updating application state...")
	s.header = fetcher.GetHeaderInfo(&s.cfg.Header)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.statusCards = fetcher.GetStatusCardData(s.cfg.StatusCards, s.cfg.GHToken)
	}()
	go func() {
		defer wg.Done()
		s.thumbFeeds = fetcher.GetThumbFeedData(s.cfg.ThumbFeeds)
	}()
	wg.Wait()
	log.Println("Application state updated.")
}

func (s *State) StartUpdateLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		s.updateState()
	}
}

func (s *State) GetHeader() *models.HeaderInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.header
}

func (s *State) GetStatusCards() []models.StatusCardSection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statusCards
}

func (s *State) GetThumbFeeds() []models.ThumbFeedSection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.thumbFeeds
}

func (s *State) GetBookmarks() []config.BookmarkCategory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bookmarks
}
