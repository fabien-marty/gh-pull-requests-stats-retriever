package gh

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/config"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/log"
	"github.com/google/go-github/v62/github"
)

type LabelEventStats struct {
	FirstAdded       *time.Time `json:"first_added"`
	LastAdded        *time.Time `json:"last_added"`
	FirstRemoved     *time.Time `json:"first_removed"`
	LastRemoved      *time.Time `json:"last_removed"`
	SecondsWithLabel int        `json:"seconds_with_label"`
	Current          bool       `json:"current"`
	Labels           []string   `json:"labels"`
}

func sortEvent(i *github.IssueEvent, j *github.IssueEvent) int {
	iTime := i.CreatedAt.GetTime()
	jTime := j.CreatedAt.GetTime()
	if iTime == nil {
		log.GetLogger().Warn("IssueEvent has no CreatedAt")
		if jTime == nil {
			return 0
		}
		return 1
	}
	if jTime == nil {
		log.GetLogger().Warn("IssueEvent has no CreatedAt")
		if iTime == nil {
			return 0
		}
		return -1
	}
	if iTime.Before(*jTime) {
		return -1
	} else if iTime.After(*jTime) {
		return 1
	} else {
		return 0
	}
}

func getStats(events []*github.IssueEvent, labelNames []string) (*LabelEventStats, error) {
	if len(labelNames) == 0 {
		return nil, errors.New("no label names provided")
	}
	res := LabelEventStats{Labels: labelNames}
	for _, event := range events {
		if event.CreatedAt == nil {
			continue
		}
		if event.Label == nil {
			continue
		}
		if !slices.Contains(labelNames, event.Label.GetName()) {
			continue
		}
		eventCreatedTime := event.CreatedAt.GetTime()
		switch event.GetEvent() {
		case "labeled":
			if res.FirstAdded == nil {
				res.FirstAdded = eventCreatedTime
			}
			res.LastAdded = eventCreatedTime
			res.Current = true
		case "unlabeled":
			if res.FirstRemoved == nil {
				res.FirstRemoved = eventCreatedTime
			}
			res.LastRemoved = eventCreatedTime
			if res.Current {
				res.Current = false
				res.SecondsWithLabel += int(eventCreatedTime.Sub(*res.LastAdded).Seconds())
			}
		default:
			return nil, fmt.Errorf("unknown event type: %s", event.GetEvent())
		}
	}
	if res.Current && res.LastAdded != nil {
		res.SecondsWithLabel += int(time.Since(*res.LastAdded).Seconds())
	}
	return &res, nil
}

func GetLabelEventStats(client *Client, config *config.Config, prNumber int) ([]LabelEventStats, error) {
	logger := log.GetLogger().With(slog.String("owner", config.Owner), slog.String("repo", config.Repo))
	options := github.ListOptions{}
	events := []*github.IssueEvent{}
	for {
		logger.Info("Fetching label issue events...", slog.Int("prNumber", prNumber), slog.Int("page", options.Page))
		evts, resp, err := client.client.Issues.ListIssueEvents(context.Background(), config.Owner, config.Repo, prNumber, nil)
		if err != nil {
			return nil, err
		}
		for _, evt := range evts {
			if evt.GetEvent() != "labeled" && evt.GetEvent() == "unlabeled" {
				continue
			}
			events = append(events, evt)
		}
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}
	slices.SortFunc(events, sortEvent)
	logger.Info("Issue label events fetched", slog.Int("count", len(events)))
	res := []LabelEventStats{}
	for _, labelNames := range config.LabelsStats.Labels {
		stats, err := getStats(events, labelNames)
		if err != nil {
			return nil, err
		}
		res = append(res, *stats)
	}
	return res, nil
}
