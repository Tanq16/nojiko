package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Header      HeaderConfig       `yaml:"header"`
	StatusCards []StatusCardConfig `yaml:"statusCards"`
	ThumbFeeds  []ThumbFeedConfig  `yaml:"thumbFeeds"`
	Bookmarks   []BookmarkCategory `yaml:"bookmarks"`
}

type HeaderConfig struct {
	Title         string `yaml:"title"`
	LogoURL       string `yaml:"logoURL"`
	ShowWeather   bool   `yaml:"showWeather"`
	WeatherAPIKey string `yaml:"weatherAPIKey"`
	City          string `yaml:"city"`
}

type StatusCardConfig struct {
	Title        string            `yaml:"title"`
	Icon         string            `yaml:"icon"`
	ShowMockData bool              `yaml:"showMockData"`
	Repositories []GitHubRepoOwner `yaml:"repositories"`
}

type GitHubRepoOwner struct {
	Owner string `yaml:"owner"`
	Repo  string `yaml:"repo"`
}

type ThumbFeedConfig struct {
	Title        string   `yaml:"title"`
	Icon         string   `yaml:"icon"`
	FeedType     string   `yaml:"feedType"`
	ShowMockData bool     `yaml:"showMockData"`
	Channels     []string `yaml:"channels"`
}

type BookmarkCategory struct {
	Category string   `yaml:"category" json:"category"`
	Color    string   `yaml:"color"    json:"color"`
	Links    []Link   `yaml:"links"    json:"links"`
	Folders  []Folder `yaml:"folders"  json:"folders"`
}

type Folder struct {
	Name  string `yaml:"name"  json:"name"`
	Icon  string `yaml:"icon"  json:"icon"`
	Links []Link `yaml:"links" json:"links"`
}

type Link struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url"  json:"url"`
	Icon string `yaml:"icon" json:"icon"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
