package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

var timeAgoRegex = regexp.MustCompile(`(\d+)\s+(minute|hour|day|week|month|year)s?\s+ago`)

func parseTimeAgo(ago string) time.Time {
	now := time.Now()
	matches := timeAgoRegex.FindStringSubmatch(strings.ToLower(ago))
	if len(matches) != 3 {
		return time.Time{}
	}

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}
	}

	unit := matches[2]
	switch unit {
	case "minute":
		return now.Add(-time.Duration(val) * time.Minute)
	case "hour":
		return now.Add(-time.Duration(val) * time.Hour)
	case "day":
		return now.AddDate(0, 0, -val)
	case "week":
		return now.AddDate(0, 0, -val*7)
	case "month":
		return now.AddDate(0, -val, 0)
	case "year":
		return now.AddDate(-val, 0, 0)
	}

	return time.Time{}
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
	if len(data.Contents.TwoColumnBrowseResultsRenderer.Tabs) < 2 {
		return nil
	}
	videoItems := data.Contents.TwoColumnBrowseResultsRenderer.Tabs[1].TabRenderer.Content.RichGridRenderer.Contents
	for _, item := range videoItems {
		if len(videos) >= 5 {
			break
		}
		vr := item.RichItemRenderer.Content.VideoRenderer
		if vr.VideoID == "" || len(vr.Title.Runs) == 0 {
			continue
		}
		var bestThumbnail string
		if len(vr.Thumbnail.Thumbnails) > 0 {
			bestThumbnail = vr.Thumbnail.Thumbnails[len(vr.Thumbnail.Thumbnails)-1].URL
		}
		videos = append(videos, YouTubeCard{
			Type:        "youtube",
			Title:       vr.Title.Runs[0].Text,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", vr.VideoID),
			Channel:     channelName,
			Published:   vr.Published.SimpleText,
			Thumbnail:   bestThumbnail,
			PublishedAt: parseTimeAgo(vr.Published.SimpleText),
		})
	}
	return videos
}

func GetThumbFeedData(configs []config.ThumbFeedConfig) []ThumbFeedSection {
	var sections []ThumbFeedSection
	scraper := newYouTubeScraper()

	for _, conf := range configs {
		section := ThumbFeedSection{
			Title:    conf.Title,
			Icon:     conf.Icon,
			FeedType: conf.FeedType,
		}
		var allVideos []YouTubeCard
		if conf.FeedType == "youtube" && len(conf.Channels) > 0 {
			for _, channel := range conf.Channels {
				fetchedVideos := scraper.getLatestVideos(channel)
				if len(fetchedVideos) > 0 {
					allVideos = append(allVideos, fetchedVideos...)
				}
				time.Sleep(200 * time.Millisecond)
			}
		}

		sort.Slice(allVideos, func(i, j int) bool {
			return allVideos[i].PublishedAt.After(allVideos[j].PublishedAt)
		})

		if conf.Limit > 0 && len(allVideos) > conf.Limit {
			allVideos = allVideos[:conf.Limit]
		}

		section.Cards = allVideos
		sections = append(sections, section)
	}
	return sections
}
