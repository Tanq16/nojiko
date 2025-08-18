package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/tanq16/nojiko/internal/config"
	"github.com/tanq16/nojiko/internal/models"
)

type Jellyfin struct {
	client  *http.Client
	config  config.StatusCardItem
	baseURL *url.URL
}

type jellyfinItemsResponse struct {
	TotalRecordCount int `json:"TotalRecordCount"`
}

func NewJellyfin(cfg config.StatusCardItem) *Jellyfin {
	bu, err := url.Parse(cfg.URL)
	if err != nil {
		log.Printf("Invalid URL for Jellyfin: %s", cfg.URL)
		return nil
	}
	return &Jellyfin{
		client:  &http.Client{Timeout: 10 * time.Second},
		config:  cfg,
		baseURL: bu,
	}
}

func (j *Jellyfin) fetchItemCount(itemType string) (int, error) {
	countURL := j.baseURL.JoinPath("/Items")
	q := countURL.Query()
	q.Set("IncludeItemTypes", itemType)
	q.Set("Recursive", "true")
	countURL.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", countURL.String(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-MediaBrowser-Token", j.config.APIKey)
	resp, err := j.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var data jellyfinItemsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	return data.TotalRecordCount, nil
}

func (j *Jellyfin) FetchStatus() (*models.ServiceStatusCard, error) {
	card := &models.ServiceStatusCard{Name: j.config.Name}
	healthURL := j.baseURL.JoinPath("/System/Info")
	req, err := http.NewRequest("GET", healthURL.String(), nil)
	if err != nil {
		card.Healthy = false
		log.Printf("Jellyfin at %s is unhealthy. Failed to create request: %v", j.config.URL, err)
		return card, nil
	}
	req.Header.Set("X-MediaBrowser-Token", j.config.APIKey)

	healthResp, err := j.client.Do(req)
	if err != nil || healthResp.StatusCode >= 400 {
		card.Healthy = false
		if err != nil {
			log.Printf("Jellyfin at %s is unhealthy. Error: %v", j.config.URL, err)
		} else {
			log.Printf("Jellyfin at %s is unhealthy. Status: %s", j.config.URL, healthResp.Status)
			healthResp.Body.Close()
		}
		return card, nil
	}
	card.Healthy = true
	healthResp.Body.Close()

	mediaTypes := map[string]string{
		"Movies": "Movie",
		"Series": "Series",
		"Songs":  "Audio",
	}
	stats := []*models.NumStat{}
	for name, apiType := range mediaTypes {
		count, err := j.fetchItemCount(apiType)
		if err != nil {
			log.Printf("Failed to fetch %s count from Jellyfin: %v", name, err)
			continue
		}
		stats = append(stats, &models.NumStat{
			Text:         name,
			DisplayValue: fmt.Sprintf("%d", count),
		})
	}
	if len(stats) > 0 {
		card.NumStats1 = stats[0]
	}
	if len(stats) > 1 {
		card.NumStats2 = stats[1]
	}
	if len(stats) > 2 {
		card.NumStats3 = stats[2]
	}
	if len(stats) > 3 {
		card.NumStats4 = stats[3]
	}
	return card, nil
}
