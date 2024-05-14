package repo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseIssueEventType(t *testing.T) {
	assert.Equal(t, IssueTimelineTypeLabeled, ParseIssueTimelineType("labeled"))
	assert.Equal(t, IssueTimelineTypeUnlabeled, ParseIssueTimelineType("unlabeled"))
	assert.Equal(t, IssueTimelineTypeReopened, ParseIssueTimelineType("reopened"))
	assert.Equal(t, IssueTimelineTypeClosed, ParseIssueTimelineType("closed"))
	assert.Equal(t, IssueTimelineTypeUnknown, ParseIssueTimelineType("foobar"))
}

func TestSortIssueEvents(t *testing.T) {
	event1 := IssueTimeline{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event2 := IssueTimeline{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      IssueTimelineTypeLabeled,
		Label:     "bar",
	}
	event2bis := IssueTimeline{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      IssueTimelineTypeLabeled,
		Label:     "bar",
	}
	event3 := IssueTimeline{
		CreatedAt: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC),
		Type:      IssueTimelineTypeLabeled,
		Label:     "bar",
	}
	events := []IssueTimeline{event3, event2bis, event1, event2}
	events = SortIssueTimelines(events)
	assert.Equal(t, []IssueTimeline{event1, event2, event2bis, event3}, events)
}
