package models

import "time"

type GitHubCard struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Stars  int    `json:"stars"`
	Issues int    `json:"issues"`
	PRs    int    `json:"prs"`
}

type NumStat struct {
	Numerator    int    `json:"numerator,omitempty"`
	Denominator  int    `json:"denominator,omitempty"`
	Unit         string `json:"unit,omitempty"`
	Text         string `json:"text"`
	DisplayValue string `json:"displayValue"`
}

type ServiceStatusCard struct {
	Name      string   `json:"name"`
	Healthy   bool     `json:"healthy"`
	NumStats1 *NumStat `json:"numStats1,omitempty"`
	NumStats2 *NumStat `json:"numStats2,omitempty"`
	NumStats3 *NumStat `json:"numStats3,omitempty"`
	NumStats4 *NumStat `json:"numStats4,omitempty"`
}

type YouTubeCard struct {
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Channel     string    `json:"channel"`
	URL         string    `json:"url"`
	Published   string    `json:"published"`
	Thumbnail   string    `json:"thumbnail"`
	PublishedAt time.Time `json:"-"`
}

type StatusCardSection struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Type  string `json:"type"`
	Cards []any  `json:"cards"`
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

type OpenMeteoResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current_weather"`
}
