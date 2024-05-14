package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
	"github.com/stretchr/testify/assert"
)

func TestGetStats1(t *testing.T) {
	event1 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event2 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "foo",
	}
	event2dup := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "foo",
	}
	event3 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event4 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeClosed,
		Label:     "",
	}
	pr := repo.PullRequest{
		Number: 1,
		State:  "closed",
	}
	stats, err := getStats([]repo.IssueTimeline{event1, event2, event2dup, event3, event4}, []string{"foo"}, pr)
	assert.NoError(t, err)
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.FirstAdded.Format(time.RFC3339))
	assert.Equal(t, "2021-01-05T00:00:00Z", stats.LastAdded.Format(time.RFC3339))
	assert.Equal(t, "2021-01-03T00:00:00Z", stats.FirstRemoved.Format(time.RFC3339))
	assert.Equal(t, "2021-01-03T00:00:00Z", stats.LastRemoved.Format(time.RFC3339))
	assert.True(t, stats.Current)
	assert.ElementsMatch(t, []string{"foo"}, stats.Labels)
	assert.Equal(t, 259200, stats.SecondsWithLabel)
}

func TestGetStats2(t *testing.T) {
	event1 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event2 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "bar",
	}
	event3 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "bar",
	}
	event4 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeClosed,
	}
	event5 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeReopened,
	}
	event6 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeMerged,
	}
	pr := repo.PullRequest{
		Number: 1,
		State:  "closed",
	}
	stats, err := getStats([]repo.IssueTimeline{event1, event2, event3, event4, event5, event6}, []string{"foo"}, pr)
	assert.NoError(t, err)
	fmt.Printf("%+v\n", stats)
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.FirstAdded.Format(time.RFC3339))
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.LastAdded.Format(time.RFC3339))
	assert.Nil(t, stats.FirstRemoved)
	assert.Nil(t, stats.LastRemoved)
	assert.True(t, stats.Current)
	assert.ElementsMatch(t, []string{"foo"}, stats.Labels)
	assert.Equal(t, 345600, stats.SecondsWithLabel)
}
