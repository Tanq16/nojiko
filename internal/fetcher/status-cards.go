package fetcher

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/tanq16/nojiko/internal/config"
	"github.com/tanq16/nojiko/internal/fetcher/services"
	"github.com/tanq16/nojiko/internal/models"
	"golang.org/x/oauth2"
)

// all service status fetchers implement this
type ServiceFetcher interface {
	FetchStatus() (*models.ServiceStatusCard, error)
}

// maps service types from the config to resp. fetchers
var serviceFetchers = map[string]func(config.StatusCardItem) ServiceFetcher{
	"adguard": func(cfg config.StatusCardItem) ServiceFetcher {
		return services.NewAdguard(cfg)
	},
}

func GetStatusCardData(configs []config.StatusCardConfig, ghToken string) []models.StatusCardSection {
	var sections []models.StatusCardSection
	ctx := context.Background()
	var client *github.Client
	if ghToken != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	for _, conf := range configs {
		section := models.StatusCardSection{
			Title: conf.Title,
			Icon:  conf.Icon,
			Type:  conf.Type,
		}
		var cards []any
		for _, item := range conf.Cards {
			var cardData any

			switch conf.Type {
			case "github":
				parts := strings.Split(item.Repo, "/")
				if len(parts) != 2 {
					log.Printf("Invalid repo format: %s. Expected 'owner/repo'", item.Repo)
					continue
				}
				owner, repoName := parts[0], parts[1]
				repo, _, err := client.Repositories.Get(ctx, owner, repoName)
				if err != nil {
					log.Printf("Error fetching repo %s: %v", item.Repo, err)
					continue
				}
				issues, _, err := client.Issues.ListByRepo(ctx, owner, repoName, &github.IssueListByRepoOptions{State: "open"})
				if err != nil {
					log.Printf("Error fetching issues for %s: %v", item.Repo, err)
					continue
				}
				openIssues := 0
				openPRs := 0
				for _, issue := range issues {
					if issue.IsPullRequest() {
						openPRs++
					} else {
						openIssues++
					}
				}
				cardData = models.GitHubCard{
					Name:   repo.GetFullName(),
					URL:    repo.GetHTMLURL(),
					Stars:  repo.GetStargazersCount(),
					Issues: openIssues,
					PRs:    openPRs,
				}

			case "service":
				if initializer, ok := serviceFetchers[item.ServiceType]; ok {
					fetcher := initializer(item)
					status, err := fetcher.FetchStatus()
					if err != nil {
						log.Printf("Error fetching status for service %s: %v", item.Name, err)
						cardData = &models.ServiceStatusCard{Name: item.Name, Healthy: false}
					} else {
						cardData = status
					}
				}
			}
			if cardData != nil {
				cards = append(cards, cardData)
			}
			time.Sleep(200 * time.Millisecond)
		}
		section.Cards = cards
		sections = append(sections, section)
	}
	return sections
}
