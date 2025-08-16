package fetcher

import "github.com/tanq16/nojiko/internal/config"

// RepoStats holds statistics for a GitHub repository.
type RepoStats struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Stars int    `json:"stars"`
	Forks int    `json:"forks"`
}

// YouTubeVideo holds information about a YouTube video.
type YouTubeVideo struct {
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

// HeaderInfo holds information for the main page header.
type HeaderInfo struct {
	Title   string       `json:"title"`
	LogoURL string       `json:"logoURL"`
	Weather *WeatherInfo `json:"weather,omitempty"`
}

// WeatherInfo holds weather data.
type WeatherInfo struct {
	TempC       int    `json:"tempC"`
	Description string `json:"description"`
}

// GetGitHubStats returns stats for configured repositories.
func GetGitHubStats(cfg config.GithubConfig) []RepoStats {
	// Placeholder implementation.
	// In a real scenario, you'd make an API call to GitHub here.
	if cfg.ShowMockData {
		return []RepoStats{
			{Name: "homelab-dashboard", URL: "#", Stars: 7300, Forks: 251},
			{Name: "dotfiles", URL: "#", Stars: 450, Forks: 32},
			{Name: "awesome-project-x", URL: "#", Stars: 1200, Forks: 109},
			{Name: "media-server-config", URL: "#", Stars: 88, Forks: 12},
			{Name: "ansible-playbooks", URL: "#", Stars: 156, Forks: 45},
			{Name: "nix-configs", URL: "#", Stars: 310, Forks: 28},
		}
	}
	return []RepoStats{}
}

// GetYouTubeVideos returns the latest videos from configured channels.
func GetYouTubeVideos(cfg config.YoutubeConfig) []YouTubeVideo {
	// Placeholder implementation.
	if cfg.ShowMockData {
		return []YouTubeVideo{
			{Title: "Making Sense of justify-content", Channel: "Kevin Powell", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+1"},
			{Title: "Zen 5 And AI Boom w/ Casey Muratori", Channel: "ThePrimeTime", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+2"},
			{Title: "Building the Lowest Rated PC", Channel: "Linus Tech Tips", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+3"},
			{Title: "Advanced Slicer settings and...", Channel: "Maker's Muse", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+4"},
			{Title: "Electronic Aircraft", Channel: "Tom Scott", URL: "#", Thumbnail: "https://placehold.co/600x400/181825/cdd6f4?text=Video+5"},
		}
	}
	return []YouTubeVideo{}
}

// GetHeaderInfo returns data for the header.
func GetHeaderInfo() *HeaderInfo {
	// Placeholder implementation.
	return &HeaderInfo{
		Title:   "Nojiko Dashboard",
		LogoURL: "logo.png",
		Weather: &WeatherInfo{
			TempC:       19,
			Description: "Partly Cloudy",
		},
	}
}
