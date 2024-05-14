package repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type repoTestAdapter struct {
	prss    map[string][]PullRequest
	eventss map[int][]IssueTimeline
}

func (r *repoTestAdapter) getKey(owner, repo string, opts GetPRsOptions) string {
	return fmt.Sprintf("%s / %s / %s / %d / %d", owner, repo, opts.State, opts.RestrictToPr, opts.UpdatedSeconds)
}

func (r *repoTestAdapter) addPRs(owner, repo string, opts GetPRsOptions, prs []PullRequest) {
	key := r.getKey(owner, repo, opts)
	if r.prss == nil {
		r.prss = make(map[string][]PullRequest)
	}
	r.prss[key] = prs
}

func (r *repoTestAdapter) addIssueEvents(prNumber int, events []IssueTimeline) {
	if r.eventss == nil {
		r.eventss = make(map[int][]IssueTimeline)
	}
	r.eventss[prNumber] = events
}

func (r *repoTestAdapter) GetPRs(owner, repo string, opts GetPRsOptions) ([]PullRequest, error) {
	key := r.getKey(owner, repo, opts)
	prs, ok := r.prss[key]
	if !ok {
		return nil, nil
	}
	return prs, nil
}

func (r *repoTestAdapter) ListIssueTimeline(owner, repo string, issueNumber int) ([]IssueTimeline, error) {
	issues, ok := r.eventss[issueNumber]
	if !ok {
		return nil, nil
	}
	return issues, nil
}

func TestGetPRs(t *testing.T) {
	adapter := repoTestAdapter{}
	service := NewService(&adapter, "owner", "repo")
	opt1 := GetPRsOptions{State: "open"}
	opt2 := GetPRsOptions{State: "all"}
	adapter.addPRs("owner", "repo", opt1, []PullRequest{{Number: 1}, {Number: 2}})
	adapter.addPRs("owner", "repo", opt2, []PullRequest{{Number: 3}, {Number: 2}})
	prs, err := service.GetPRs([]GetPRsOptions{opt1, opt2})
	assert.NoError(t, err)
	assert.Equal(t, []PullRequest{{Number: 1}, {Number: 2}, {Number: 3}}, prs)
}

func TestListIssueTimeline(t *testing.T) {
	adapter := repoTestAdapter{}
	event1 := IssueTimeline{}
	event2 := IssueTimeline{}
	adapter.addIssueEvents(1, []IssueTimeline{event1, event2})
	service := NewService(&adapter, "owner", "repo")
	events, err := service.ListIssueTimeline(1)
	assert.NoError(t, err)
	assert.Equal(t, []IssueTimeline{event1, event2}, events)
}
