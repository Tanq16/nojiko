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

type Adguard struct {
	client  *http.Client
	config  config.StatusCardItem
	baseURL *url.URL
}

type adguardStatsResponse struct {
	NumDNSQueries       int     `json:"num_dns_queries"`
	NumBlockedFiltering int     `json:"num_blocked_filtering"`
	AvgProcessingTime   float64 `json:"avg_processing_time"`
}

func NewAdguard(cfg config.StatusCardItem) *Adguard {
	bu, err := url.Parse(cfg.URL)
	if err != nil {
		log.Printf("Invalid URL for Adguard: %s", cfg.URL)
		return nil
	}
	return &Adguard{
		client:  &http.Client{Timeout: 5 * time.Second},
		config:  cfg,
		baseURL: bu,
	}
}

func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%.1fm", float64(n)/1000000.0)
}

func (a *Adguard) FetchStatus() (*models.ServiceStatusCard, error) {
	card := &models.ServiceStatusCard{Name: a.config.Name}
	healthResp, err := a.client.Get(a.config.URL)
	if err != nil || healthResp.StatusCode >= 400 {
		card.Healthy = false
		if err != nil {
			log.Printf("Adguard at %s is unhealthy. Error: %v", a.config.URL, err)
		} else {
			log.Printf("Adguard at %s is unhealthy. Status: %s", a.config.URL, healthResp.Status)
			healthResp.Body.Close()
		}
		return card, nil
	}
	card.Healthy = true
	healthResp.Body.Close()

	statsURL := a.baseURL.JoinPath("/control/stats")
	req, err := http.NewRequest("GET", statsURL.String(), nil)
	if err != nil {
		log.Printf("Failed to create stats request for Adguard: %v", err)
		return card, nil
	}
	req.SetBasicAuth(a.config.Username, a.config.Password)
	statsResp, err := a.client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch stats from Adguard: %v", err)
		card.Healthy = false
		return card, nil
	}
	defer statsResp.Body.Close()

	if statsResp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch Adguard stats, status: %s", statsResp.Status)
		card.Healthy = false
		return card, nil
	}
	body, err := io.ReadAll(statsResp.Body)
	if err != nil {
		log.Printf("Failed to read Adguard stats response body: %v", err)
		return card, nil
	}
	var stats adguardStatsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		log.Printf("Failed to unmarshal Adguard stats: %v", err)
		return card, nil
	}
	blockedRatio := 0.0
	if stats.NumDNSQueries > 0 {
		blockedRatio = (float64(stats.NumBlockedFiltering) / float64(stats.NumDNSQueries)) * 100
	}
	card.NumStats1 = &models.NumStat{
		Text:         "Queries",
		DisplayValue: formatNumber(stats.NumDNSQueries),
	}
	card.NumStats2 = &models.NumStat{
		Text:         "Blocked",
		DisplayValue: formatNumber(stats.NumBlockedFiltering),
	}
	card.NumStats3 = &models.NumStat{
		Text:         "Ratio",
		DisplayValue: fmt.Sprintf("%.1f", blockedRatio),
		Unit:         "%",
	}
	return card, nil
}
