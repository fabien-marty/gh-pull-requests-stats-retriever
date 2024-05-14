package repo

type GetPRsOptions struct {
	State          string
	RestrictToPr   int
	UpdatedSeconds int
}

type Port interface {
	GetPRs(owner, repo string, opts GetPRsOptions) ([]PullRequest, error)
	ListIssueTimeline(owner, repo string, issueNumber int) ([]IssueTimeline, error)
}
