package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	GHToken     string             `yaml:"ghToken"`
	Header      HeaderConfig       `yaml:"header"`
	StatusCards []StatusCardConfig `yaml:"statusCards"`
	ThumbFeeds  []ThumbFeedConfig  `yaml:"thumbFeeds"`
	Bookmarks   []BookmarkCategory `yaml:"bookmarks"`
}

type HeaderConfig struct {
	Title       string  `yaml:"title"`
	LogoURL     string  `yaml:"logoURL"`
	ShowLogo    bool    `yaml:"showLogo"`
	ShowWeather bool    `yaml:"showWeather"`
	Latitude    float64 `yaml:"latitude"`
	Longitude   float64 `yaml:"longitude"`
}

type StatusCardConfig struct {
	Title string           `yaml:"title"`
	Icon  string           `yaml:"icon"`
	Cards []StatusCardItem `yaml:"cards"`
}

type StatusCardItem struct {
	Type   string `yaml:"type"`
	Owner  string `yaml:"owner,omitempty"`
	Repo   string `yaml:"repo,omitempty"`
	Name   string `yaml:"name,omitempty"`
	APIKey string `yaml:"apikey,omitempty"`
	URL    string `yaml:"url,omitempty"`
}

type ThumbFeedConfig struct {
	Title    string   `yaml:"title"`
	Icon     string   `yaml:"icon"`
	FeedType string   `yaml:"feedType"`
	Channels []string `yaml:"channels"`
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

	for i := range cfg.StatusCards {
		if cfg.StatusCards[i].Icon == "" {
			cfg.StatusCards[i].Icon = "bar-chart-3"
		}
	}
	for i := range cfg.ThumbFeeds {
		if cfg.ThumbFeeds[i].Icon == "" {
			cfg.ThumbFeeds[i].Icon = "tv-minimal"
		}
	}

	for i := range cfg.Bookmarks {
		for j := range cfg.Bookmarks[i].Links {
			if cfg.Bookmarks[i].Links[j].Icon == "" {
				cfg.Bookmarks[i].Links[j].Icon = "default"
			}
		}
		for j := range cfg.Bookmarks[i].Folders {
			if cfg.Bookmarks[i].Folders[j].Icon == "" {
				cfg.Bookmarks[i].Folders[j].Icon = "folder"
			}
			for k := range cfg.Bookmarks[i].Folders[j].Links {
				if cfg.Bookmarks[i].Folders[j].Links[k].Icon == "" {
					cfg.Bookmarks[i].Folders[j].Links[k].Icon = "default"
				}
			}
		}
	}
	return &cfg, nil
}
