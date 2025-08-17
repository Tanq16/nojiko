package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/tanq16/nojiko/internal/config"
)

type youtubeScraper struct {
	client *http.Client
}

func newYouTubeScraper() *youtubeScraper {
	return &youtubeScraper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Simplified structs to navigate the complex ytInitialData JSON
type ytInitialData struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						RichGridRenderer struct {
							Contents []struct {
								RichItemRenderer struct {
									Content struct {
										VideoRenderer videoRenderer `json:"videoRenderer"`
									} `json:"content"`
								} `json:"richItemRenderer"`
							} `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type videoRenderer struct {
	VideoID   string        `json:"videoId"`
	Title     runs          `json:"title"`
	Published simpleText    `json:"publishedTimeText"`
	Thumbnail thumbnailData `json:"thumbnail"`
}

type runs struct {
	Runs []struct {
		Text string `json:"text"`
	} `json:"runs"`
}

type simpleText struct {
	SimpleText string `json:"simpleText"`
}

type thumbnailData struct {
	Thumbnails []struct {
		URL   string `json:"url"`
		Width int    `json:"width"`
	} `json:"thumbnails"`
}

func (s *youtubeScraper) getLatestVideos(channelName string) []YouTubeCard {
	url := fmt.Sprintf("https://www.youtube.com/@%s/videos", channelName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for channel %s: %v", channelName, err)
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Error fetching page for channel %s: %v", channelName, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching page for channel %s: status code %d", channelName, resp.StatusCode)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body for channel %s: %v", channelName, err)
		return nil
	}
	re := regexp.MustCompile(`var ytInitialData = ({.*?});`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		log.Printf("Could not find video data for channel '%s'", channelName)
		return nil
	}
	var data ytInitialData
	if err := json.Unmarshal(matches[1], &data); err != nil {
		log.Printf("Failed to parse video data for '%s'. Error: %v", channelName, err)
		return nil
	}

	var videos []YouTubeCard
	// The video data is typically in the second tab
	if len(data.Contents.TwoColumnBrowseResultsRenderer.Tabs) < 2 {
		return nil
	}
	videoItems := data.Contents.TwoColumnBrowseResultsRenderer.Tabs[1].TabRenderer.Content.RichGridRenderer.Contents
	for _, item := range videoItems {
		if len(videos) >= 5 { // Limit to 5 videos
			break
		}
		vr := item.RichItemRenderer.Content.VideoRenderer
		if vr.VideoID == "" || len(vr.Title.Runs) == 0 {
			continue
		}
		// Get the highest quality thumbnail
		var bestThumbnail string
		if len(vr.Thumbnail.Thumbnails) > 0 {
			bestThumbnail = vr.Thumbnail.Thumbnails[len(vr.Thumbnail.Thumbnails)-1].URL
		}
		videos = append(videos, YouTubeCard{
			Type:      "youtube",
			Title:     vr.Title.Runs[0].Text,
			URL:       fmt.Sprintf("https://www.youtube.com/watch?v=%s", vr.VideoID),
			Channel:   channelName,
			Published: vr.Published.SimpleText,
			Thumbnail: bestThumbnail,
		})
	}
	return videos
}

func GetThumbFeedData(configs []config.ThumbFeedConfig) []ThumbFeedSection {
	var sections []ThumbFeedSection
	scraper := newYouTubeScraper()
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, conf := range configs {
		wg.Add(1)
		go func(conf config.ThumbFeedConfig) {
			defer wg.Done()
			section := ThumbFeedSection{
				Title:    conf.Title,
				Icon:     conf.Icon,
				FeedType: conf.FeedType,
			}
			var cards []YouTubeCard
			if conf.FeedType == "youtube" && len(conf.Channels) > 0 {
				var cardWg sync.WaitGroup
				var cardMu sync.Mutex
				for _, channel := range conf.Channels {
					cardWg.Add(1)
					go func(channelName string) {
						defer cardWg.Done()
						time.Sleep(1 * time.Second)
						fetchedVideos := scraper.getLatestVideos(channelName)
						if len(fetchedVideos) > 0 {
							cardMu.Lock()
							cards = append(cards, fetchedVideos...)
							cardMu.Unlock()
						}
					}(channel)
				}
				cardWg.Wait()
			}
			section.Cards = cards
			mu.Lock()
			sections = append(sections, section)
			mu.Unlock()
		}(conf)
	}
	wg.Wait()
	return sections
}
