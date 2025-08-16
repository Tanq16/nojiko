package fetcher

import (
	"context"
	"log"
	"sync"

	"github.com/google/go-github/v74/github"
	"github.com/tanq16/nojiko/internal/config"
	"golang.org/x/oauth2"
)

func GetStatusCardData(configs []config.StatusCardConfig, ghToken string) []StatusCardSection {
	var sections []StatusCardSection
	var wg sync.WaitGroup
	var mu sync.Mutex

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
		wg.Add(1)
		go func(conf config.StatusCardConfig) {
			defer wg.Done()
			section := StatusCardSection{
				Title: conf.Title,
				Icon:  conf.Icon,
			}
			var cards []any
			var cardWg sync.WaitGroup
			var cardMu sync.Mutex

			for _, item := range conf.Cards {
				cardWg.Add(1)
				go func(item config.StatusCardItem) {
					defer cardWg.Done()
					switch item.Type {
					case "github":
						repo, _, err := client.Repositories.Get(ctx, item.Owner, item.Repo)
						if err != nil {
							log.Printf("Error fetching repo %s/%s: %v", item.Owner, item.Repo, err)
							return
						}
						issues, _, err := client.Issues.ListByRepo(ctx, item.Owner, item.Repo, &github.IssueListByRepoOptions{State: "open"})
						if err != nil {
							log.Printf("Error fetching issues for %s/%s: %v", item.Owner, item.Repo, err)
							return
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
						card := GitHubCard{
							Type:   "github",
							Name:   repo.GetFullName(),
							URL:    repo.GetHTMLURL(),
							Stars:  repo.GetStargazersCount(),
							Issues: openIssues,
							PRs:    openPRs,
						}
						cardMu.Lock()
						cards = append(cards, card)
						cardMu.Unlock()

					case "service":
						card := ServiceStatusCard{
							Type: "service",
							Name: item.Name,
						}
						cardMu.Lock()
						cards = append(cards, card)
						cardMu.Unlock()
					}
				}(item)
			}
			cardWg.Wait()
			section.Cards = cards
			mu.Lock()
			sections = append(sections, section)
			mu.Unlock()
		}(conf)
	}
	wg.Wait()
	return sections
}
