package repo

import (
	"slices"
	"time"
)

type IssueTimelineType string

const IssueTimelineTypeLabeled IssueTimelineType = "labeled"
const IssueTimelineTypeUnlabeled IssueTimelineType = "unlabeled"
const IssueTimelineTypeClosed IssueTimelineType = "closed"
const IssueTimelineTypeMerged IssueTimelineType = "merged"
const IssueTimelineTypeReopened IssueTimelineType = "reopened"
const IssueTimelineTypeUnknown IssueTimelineType = "unknown"

type IssueTimeline struct {
	Type      IssueTimelineType
	CreatedAt time.Time
	Label     string
}

func sortIssueTimeline(i IssueTimeline, j IssueTimeline) int {
	if i.CreatedAt.Before(j.CreatedAt) {
		return -1
	} else if i.CreatedAt.After(j.CreatedAt) {
		return 1
	} else {
		return 0
	}
}

func SortIssueTimelines(events []IssueTimeline) []IssueTimeline {
	slices.SortFunc(events, sortIssueTimeline)
	return events
}

func ParseIssueTimelineType(typ string) IssueTimelineType {
	switch typ {
	case "labeled":
		return IssueTimelineTypeLabeled
	case "unlabeled":
		return IssueTimelineTypeUnlabeled
	case "reopened":
		return IssueTimelineTypeReopened
	case "closed":
		return IssueTimelineTypeClosed
	default:
		return IssueTimelineTypeUnknown
	}
}
