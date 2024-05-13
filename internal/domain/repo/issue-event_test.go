package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseIssueEventType(t *testing.T) {
	assert.Equal(t, IssueEventTypeLabeled, ParseIssueEventType("labeled"))
	assert.Equal(t, IssueEventTypeUnlabeled, ParseIssueEventType("unlabeled"))
	assert.Equal(t, IssueEventTypeReopened, ParseIssueEventType("reopened"))
	assert.Equal(t, IssueEventTypeClosed, ParseIssueEventType("closed"))
	assert.Equal(t, IssueEventTypeUnknown, ParseIssueEventType("foobar"))
}

func TestSortIssueEvents(t *testing.T) {
	event1 := IssueEvent{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      IssueEventTypeLabeled,
		Label:     "foo",
	}
	event2 := IssueEvent{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      IssueEventTypeLabeled,
		Label:     "bar",
	}
	event2bis := IssueEvent{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      IssueEventTypeLabeled,
		Label:     "bar",
	}
	event3 := IssueEvent{
		CreatedAt: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC),
		Type:      IssueEventTypeLabeled,
		Label:     "bar",
	}
	events := []IssueEvent{event3, event2bis, event1, event2}
	events = SortIssueEvents(events)
	assert.Equal(t, []IssueEvent{event1, event2, event2bis, event3}, events)
}
