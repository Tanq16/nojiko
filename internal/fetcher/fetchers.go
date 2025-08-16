package fetcher

import (
	"fmt"

	"github.com/tanq16/nojiko/internal/config"
)

// Generic card data structures
type GitHubCard struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Stars int    `json:"stars"`
	Forks int    `json:"forks"`
}

type YouTubeCard struct {
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

// Section data structures
type StatusCardSection struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Cards any    `json:"cards"`
}

type ThumbFeedSection struct {
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	FeedType string `json:"feedType"`
	Cards    any    `json:"cards"`
}

type HeaderInfo struct {
	Title    string       `json:"title"`
	LogoURL  string       `json:"logoURL"`
	ShowLogo bool         `json:"showLogo"`
	Weather  *WeatherInfo `json:"weather,omitempty"`
}

type WeatherInfo struct {
	TempC       float64 `json:"tempC"`
	Description string  `json:"description"`
}

// OpenMeteo specific response structure
type OpenMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current_weather"`
}

func GetStatusCardData(configs []config.StatusCardConfig) []StatusCardSection {
	var sections []StatusCardSection
	for _, conf := range configs {
		section := StatusCardSection{
			Title: conf.Title,
			Icon:  conf.Icon,
		}
		if conf.ShowMockData {
			var cards []GitHubCard
			for _, repo := range conf.Repositories {
				cards = append(cards, GitHubCard{
					Name:  fmt.Sprintf("%s/%s", repo.Owner, repo.Repo),
					URL:   "#",
					Stars: 100,
					Forks: 20,
				})
			}
			section.Cards = cards
		}
		sections = append(sections, section)
	}
	return sections
}

func GetThumbFeedData(configs []config.ThumbFeedConfig) []ThumbFeedSection {
	var sections []ThumbFeedSection
	for _, conf := range configs {
		section := ThumbFeedSection{
			Title:    conf.Title,
			Icon:     conf.Icon,
			FeedType: conf.FeedType,
		}
		if conf.ShowMockData {
			if conf.FeedType == "youtube" {
				section.Cards = []YouTubeCard{
					{Title: "Making Sense of justify-content", Channel: "Kevin Powell", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+1"},
					{Title: "Zen 5 And AI Boom", Channel: "ThePrimeTime", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+2"},
				}
			}
		}
		sections = append(sections, section)
	}
	return sections
}
