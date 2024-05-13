package stats

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
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

func getStats(events []repo.IssueEvent, labelNames []string, pr repo.PullRequest) (*LabelEventStats, error) {
	var opened bool = true
	if len(labelNames) == 0 {
		return nil, errors.New("no label names provided")
	}
	res := LabelEventStats{Labels: labelNames}
	for _, event := range events {
		if !slices.Contains(labelNames, event.Label) {
			continue
		}
		eventCreatedTime := event.CreatedAt
		switch event.Type {
		case repo.IssueEventTypeLabeled:
			if res.FirstAdded == nil {
				res.FirstAdded = &eventCreatedTime
			}
			res.LastAdded = &eventCreatedTime
			res.Current = true
		case repo.IssueEventTypeUnlabeled:
			if res.FirstRemoved == nil {
				res.FirstRemoved = &eventCreatedTime
			}
			res.LastRemoved = &eventCreatedTime
			if res.Current {
				res.Current = false
				res.SecondsWithLabel += int(eventCreatedTime.Sub(*res.LastAdded).Seconds())
			}
		case repo.IssueEventTypeClosed, repo.IssueEventTypeMerged:
			if res.Current {
				res.SecondsWithLabel += int(eventCreatedTime.Sub(*res.LastAdded).Seconds())
			}
			opened = false
		case repo.IssueEventTypeReopened:
			opened = true
		default:
			return nil, fmt.Errorf("unknown event type: %s", event.Type)
		}
	}
	if res.Current && res.LastAdded != nil && pr.State == "open" && opened {
		res.SecondsWithLabel += int(time.Since(*res.LastAdded).Seconds())
	}
	return &res, nil
}

func ComputeLabelEventStats(repoService *repo.Service, pr repo.PullRequest, labels [][]string) ([]LabelEventStats, error) {
	events, err := repoService.ListIssueEvents(pr.Number)
	if err != nil {
		return nil, err
	}
	events = repo.SortIssueEvents(events)
	res := []LabelEventStats{}
	for _, labelNames := range labels {
		stats, err := getStats(events, labelNames, pr)
		if err != nil {
			return nil, err
		}
		res = append(res, *stats)
	}
	return res, nil
}
