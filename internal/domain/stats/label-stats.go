package stats

import (
	"errors"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
)

type LabelEventStats struct {
	LabelsGroup           []string   `json:"labels_group"`
	FirstFlagged          *time.Time `json:"first_flagged"`
	LastFlagged           *time.Time `json:"last_flagged"`
	lastFlaggedOrReopened *time.Time
	FirstUnflagged        *time.Time `json:"first_unflagged"`
	LastUnflagged         *time.Time `json:"last_unflagged"`
	SecondsWithFlag       int        `json:"seconds_with_flag"`
	labelGroup            labelGroup
	CurrentFlag           bool `json:"current_flag"`
}

func switchState(stats *LabelEventStats, eventCreatedTime time.Time) {
	if stats.labelGroup.isFlagged() {
		// we switch to flagged state
		if stats.FirstFlagged == nil {
			stats.FirstFlagged = &eventCreatedTime
		}
		stats.LastFlagged = &eventCreatedTime
		stats.lastFlaggedOrReopened = &eventCreatedTime
		stats.CurrentFlag = true
	} else {
		// we switch to unflagged state
		if stats.FirstUnflagged == nil {
			stats.FirstUnflagged = &eventCreatedTime
		}
		stats.LastUnflagged = &eventCreatedTime
		stats.SecondsWithFlag += int(eventCreatedTime.Sub(*stats.lastFlaggedOrReopened).Seconds())
		stats.CurrentFlag = false
	}
}

func getStats(events []repo.IssueTimeline, labelNames []string, pr repo.PullRequest) (*LabelEventStats, error) {
	var opened bool = true
	if len(labelNames) == 0 {
		return nil, errors.New("no label names provided")
	}
	res := LabelEventStats{labelGroup: newLabelGroup(labelNames), LabelsGroup: labelNames}
	res.CurrentFlag = res.labelGroup.isFlagged()
	if res.CurrentFlag {
		// we start with a flagged state
		res.FirstFlagged = pr.CreatedAt
		res.LastFlagged = pr.CreatedAt
		res.lastFlaggedOrReopened = pr.CreatedAt
	}
	for _, event := range events {
		eventCreatedTime := event.CreatedAt
		switch event.Type {
		case repo.IssueTimelineTypeLabeled:
			if !res.labelGroup.label(event.Label) {
				// tag not related to this label group
				continue
			}
			if res.CurrentFlag != res.labelGroup.isFlagged() {
				switchState(&res, eventCreatedTime)
			}
		case repo.IssueTimelineTypeUnlabeled:
			if !res.labelGroup.unlabel(event.Label) {
				// tag not related to this label group
				continue
			}
			if res.CurrentFlag != res.labelGroup.isFlagged() {
				switchState(&res, eventCreatedTime)
			}
		case repo.IssueTimelineTypeClosed, repo.IssueTimelineTypeMerged:
			if res.CurrentFlag {
				res.SecondsWithFlag += int(eventCreatedTime.Sub(*res.lastFlaggedOrReopened).Seconds())
			}
			opened = false
		case repo.IssueTimelineTypeReopened:
			opened = true
			res.lastFlaggedOrReopened = &eventCreatedTime
		}
	}
	if res.CurrentFlag && res.LastFlagged != nil && pr.State == "open" && opened {
		res.SecondsWithFlag += int(time.Since(*res.lastFlaggedOrReopened).Seconds())
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
