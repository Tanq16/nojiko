package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/tanq16/nojiko/internal/config"
)

func GetHeaderInfo(cfg *config.HeaderConfig) *HeaderInfo {
	logoURL := "logo.png" // Default logo
	if cfg.LogoURL != "" {
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Head(cfg.LogoURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			logoURL = cfg.LogoURL
		}
		if err != nil {
			log.Printf("Could not verify logo URL %s: %v. Falling back to default.", cfg.LogoURL, err)
		}
	}
	info := &HeaderInfo{
		Title:    cfg.Title,
		LogoURL:  logoURL,
		ShowLogo: cfg.ShowLogo,
	}
	if cfg.ShowWeather {
		info.Weather = getOpenMeteoWeather(cfg.Latitude, cfg.Longitude)
	}
	return info
}

func getOpenMeteoWeather(lat, lon float64) *WeatherInfo {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch weather data: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch weather data, status code: %d", resp.StatusCode)
		return nil
	}
	var weatherData OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		log.Printf("Failed to decode weather data: %v", err)
		return nil
	}
	return &WeatherInfo{
		TempC:       math.Round(weatherData.CurrentWeather.Temperature),
		Description: weatherCodeToDescription(weatherData.CurrentWeather.WeatherCode),
	}
}

func weatherCodeToDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1, 2, 3:
		return "Mainly clear"
	case 45, 48:
		return "Fog"
	case 51, 53, 55:
		return "Drizzle"
	case 56, 57:
		return "Freezing Drizzle"
	case 61, 63, 65:
		return "Rain"
	case 66, 67:
		return "Freezing Rain"
	case 71, 73, 75:
		return "Snow fall"
	case 77:
		return "Snow grains"
	case 80, 81, 82:
		return "Rain showers"
	case 85, 86:
		return "Snow showers"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}
