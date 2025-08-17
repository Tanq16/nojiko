package fetcher

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v74/github"
	"github.com/tanq16/nojiko/internal/config"
	"golang.org/x/oauth2"
)

func GetStatusCardData(configs []config.StatusCardConfig, ghToken string) []StatusCardSection {
	var sections []StatusCardSection

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
		section := StatusCardSection{
			Title: conf.Title,
			Icon:  conf.Icon,
			Type:  conf.Type,
		}
		var cards []any

		for _, item := range conf.Cards {
			var cardData any

			switch conf.Type {
			case "github":
				repo, _, err := client.Repositories.Get(ctx, item.Owner, item.Repo)
				if err != nil {
					log.Printf("Error fetching repo %s/%s: %v", item.Owner, item.Repo, err)
					continue
				}
				issues, _, err := client.Issues.ListByRepo(ctx, item.Owner, item.Repo, &github.IssueListByRepoOptions{State: "open"})
				if err != nil {
					log.Printf("Error fetching issues for %s/%s: %v", item.Owner, item.Repo, err)
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
				cardData = GitHubCard{
					Name:   repo.GetFullName(),
					URL:    repo.GetHTMLURL(),
					Stars:  repo.GetStargazersCount(),
					Issues: openIssues,
					PRs:    openPRs,
				}

			case "service":
				cardData = ServiceStatusCard{
					Name: item.Name,
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
