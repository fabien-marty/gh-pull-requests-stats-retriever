package stats

import (
	"testing"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
	"github.com/stretchr/testify/assert"
)

func TestGetStatsSingle1(t *testing.T) {
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
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.FirstFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-05T00:00:00Z", stats.LastFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-03T00:00:00Z", stats.FirstUnflagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-03T00:00:00Z", stats.LastUnflagged.Format(time.RFC3339))
	assert.True(t, stats.CurrentFlag)
	assert.Equal(t, 259200, stats.SecondsWithFlag)
}

func TestGetStatsSingle2(t *testing.T) {
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
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.FirstFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.LastFlagged.Format(time.RFC3339))
	assert.Nil(t, stats.FirstUnflagged)
	assert.Nil(t, stats.LastUnflagged)
	assert.True(t, stats.CurrentFlag)
	assert.Equal(t, 345600, stats.SecondsWithFlag)
}

func TestGetStatsMultiple(t *testing.T) {
	event1 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event2 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "bar",
	}
	event3 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "foo",
	}
	event4 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event5 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 16, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeClosed,
		Label:     "",
	}
	pr := repo.PullRequest{
		Number: 1,
		State:  "closed",
	}
	stats, err := getStats([]repo.IssueTimeline{event1, event2, event3, event4, event5}, []string{"foo", "bar"}, pr)
	assert.NoError(t, err)
	assert.Equal(t, "2021-01-02T00:00:00Z", stats.FirstFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-15T00:00:00Z", stats.LastFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-06T00:00:00Z", stats.FirstUnflagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-06T00:00:00Z", stats.LastUnflagged.Format(time.RFC3339))
	assert.True(t, stats.CurrentFlag)
	assert.Equal(t, 432000, stats.SecondsWithFlag)
}

func TestGetStatsNegative(t *testing.T) {
	event1 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeLabeled,
		Label:     "foo",
	}
	event2 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeUnlabeled,
		Label:     "foo",
	}
	event3 := repo.IssueTimeline{
		CreatedAt: time.Date(2021, 1, 8, 0, 0, 0, 0, time.UTC),
		Type:      repo.IssueTimelineTypeMerged,
		Label:     "",
	}
	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	pr := repo.PullRequest{
		Number:    1,
		State:     "closed",
		CreatedAt: &start,
	}
	stats, err := getStats([]repo.IssueTimeline{event1, event2, event3}, []string{"-foo"}, pr)
	assert.NoError(t, err)
	assert.Equal(t, "2021-01-01T00:00:00Z", stats.FirstFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-04T00:00:00Z", stats.LastFlagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-02T00:00:00Z", stats.FirstUnflagged.Format(time.RFC3339))
	assert.Equal(t, "2021-01-02T00:00:00Z", stats.LastUnflagged.Format(time.RFC3339))
	assert.True(t, stats.CurrentFlag)
	assert.Equal(t, 432000, stats.SecondsWithFlag)
}
