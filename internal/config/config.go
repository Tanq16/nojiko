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
	Type  string           `yaml:"type"`
	Cards []StatusCardItem `yaml:"cards"`
}

type StatusCardItem struct {
	Repo        string `yaml:"repo,omitempty"`
	Name        string `yaml:"name,omitempty"`
	ServiceType string `yaml:"serviceType,omitempty"`
	URL         string `yaml:"url,omitempty"`
	APIKey      string `yaml:"apiKey,omitempty"`
	Username    string `yaml:"username,omitempty"`
	Password    string `yaml:"password,omitempty"`
}

type ThumbFeedConfig struct {
	Title    string   `yaml:"title"`
	Icon     string   `yaml:"icon"`
	FeedType string   `yaml:"feedType"`
	Channels []string `yaml:"channels"`
	Limit    int      `yaml:"limit,omitempty"`
}

type BookmarkCategory struct {
	Category string   `yaml:"category" json:"category"`
	Color    string   `yaml:"color"    json:"color"`
	Folded   bool     `yaml:"folded,omitempty" json:"folded"`
	Links    []Link   `yaml:"links"    json:"links"`
	Folders  []Folder `yaml:"folders"  json:"folders"`
}

type Folder struct {
	Name   string `yaml:"name"  json:"name"`
	Icon   string `yaml:"icon"  json:"icon"`
	Folded bool   `yaml:"folded,omitempty" json:"folded"`
	Links  []Link `yaml:"links" json:"links"`
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
