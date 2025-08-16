package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

// Config is the top-level configuration struct.
type Config struct {
	Server    ServerConfig       `yaml:"server"`
	Github    GithubConfig       `yaml:"github"`
	Youtube   YoutubeConfig      `yaml:"youtube"`
	Bookmarks []BookmarkCategory `yaml:"bookmarks"`
	Services  []Service          `yaml:"services"`
	Header    HeaderConfig       `yaml:"header"`
}

// ServerConfig holds server-related settings.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// GithubConfig holds GitHub-related settings.
type GithubConfig struct {
	Username     string   `yaml:"username"`
	Repositories []string `yaml:"repositories"`
	ShowMockData bool     `yaml:"showMockData"`
}

// YoutubeConfig holds YouTube-related settings.
type YoutubeConfig struct {
	APIKey       string   `yaml:"apiKey"`
	Channels     []string `yaml:"channels"`
	ShowMockData bool     `yaml:"showMockData"`
}

// HeaderConfig holds settings for the main header.
type HeaderConfig struct {
	Title         string `yaml:"title"`
	LogoURL       string `yaml:"logoURL"`
	ShowWeather   bool   `yaml:"showWeather"`
	WeatherAPIKey string `yaml:"weatherAPIKey"`
	City          string `yaml:"city"`
}

// BookmarkCategory represents a category of bookmarks.
type BookmarkCategory struct {
	Category string   `yaml:"category" json:"category"`
	Color    string   `yaml:"color"    json:"color"`
	Links    []Link   `yaml:"links"    json:"links"`
	Folders  []Folder `yaml:"folders"  json:"folders"`
}

// Folder represents a collapsible folder of links.
type Folder struct {
	Name  string `yaml:"name"  json:"name"`
	Icon  string `yaml:"icon"  json:"icon"`
	Links []Link `yaml:"links" json:"links"`
}

// Link represents a single bookmark.
type Link struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url"  json:"url"`
	Icon string `yaml:"icon" json:"icon"`
}

// Service represents a self-hosted service to be monitored.
type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// Load reads and parses the configuration file.
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
