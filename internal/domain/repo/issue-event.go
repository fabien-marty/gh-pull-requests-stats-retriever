package repo

import (
	"slices"
	"time"
)

type IssueEventType string

const IssueEventTypeLabeled IssueEventType = "labeled"
const IssueEventTypeUnlabeled IssueEventType = "unlabeled"
const IssueEventTypeClosed IssueEventType = "closed"
const IssueEventTypeMerged IssueEventType = "merged"
const IssueEventTypeReopened IssueEventType = "reopened"
const IssueEventTypeUnknown IssueEventType = "unknown"

type IssueEvent struct {
	Type      IssueEventType
	CreatedAt time.Time
	Label     string
}

func sortIssueEvent(i IssueEvent, j IssueEvent) int {
	if i.CreatedAt.Before(j.CreatedAt) {
		return -1
	} else if i.CreatedAt.After(j.CreatedAt) {
		return 1
	} else {
		return 0
	}
}

func SortIssueEvents(events []IssueEvent) []IssueEvent {
	slices.SortFunc(events, sortIssueEvent)
	return events
}

func ParseIssueEventType(typ string) IssueEventType {
	switch typ {
	case "labeled":
		return IssueEventTypeLabeled
	case "unlabeled":
		return IssueEventTypeUnlabeled
	case "reopened":
		return IssueEventTypeReopened
	case "closed":
		return IssueEventTypeClosed
	default:
		return IssueEventTypeUnknown
	}
}
