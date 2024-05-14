package stats

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
)

type LabelEventStats struct {
	FirstAdded          *time.Time `json:"first_added"`
	LastAdded           *time.Time `json:"last_added"`
	lastAddedOrReopened *time.Time
	FirstRemoved        *time.Time `json:"first_removed"`
	LastRemoved         *time.Time `json:"last_removed"`
	SecondsWithLabel    int        `json:"seconds_with_label"`
	Current             bool       `json:"current"`
	Labels              []string   `json:"labels"`
}

func getStats(events []repo.IssueTimeline, labelNames []string, pr repo.PullRequest) (*LabelEventStats, error) {
	var opened bool = true
	if len(labelNames) == 0 {
		return nil, errors.New("no label names provided")
	}
	res := LabelEventStats{Labels: labelNames}
	for _, event := range events {
		eventCreatedTime := event.CreatedAt
		switch event.Type {
		case repo.IssueTimelineTypeLabeled:
			if !slices.Contains(labelNames, event.Label) {
				continue
			}
			if res.FirstAdded == nil {
				res.FirstAdded = &eventCreatedTime
			}
			res.LastAdded = &eventCreatedTime
			res.lastAddedOrReopened = &eventCreatedTime
			res.Current = true
		case repo.IssueTimelineTypeUnlabeled:
			if !slices.Contains(labelNames, event.Label) {
				continue
			}
			if res.FirstRemoved == nil {
				res.FirstRemoved = &eventCreatedTime
			}
			res.LastRemoved = &eventCreatedTime
			if res.Current {
				res.Current = false
				res.SecondsWithLabel += int(eventCreatedTime.Sub(*res.lastAddedOrReopened).Seconds())
			}
		case repo.IssueTimelineTypeClosed, repo.IssueTimelineTypeMerged:
			if res.Current {
				res.SecondsWithLabel += int(eventCreatedTime.Sub(*res.lastAddedOrReopened).Seconds())
			}
			opened = false
		case repo.IssueTimelineTypeReopened:
			opened = true
			res.lastAddedOrReopened = &eventCreatedTime
		default:
			return nil, fmt.Errorf("unknown event type: %s", event.Type)
		}
	}
	if res.Current && res.LastAdded != nil && pr.State == "open" && opened {
		res.SecondsWithLabel += int(time.Since(*res.lastAddedOrReopened).Seconds())
	}
	return &res, nil
}

func ComputeLabelEventStats(repoService *repo.Service, pr repo.PullRequest, labels [][]string) ([]LabelEventStats, error) {
	events, err := repoService.ListIssueTimeline(pr.Number)
	if err != nil {
		return nil, err
	}
	events = repo.SortIssueTimelines(events)
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
