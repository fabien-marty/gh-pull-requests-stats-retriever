package repo

type GetPRsOptions struct {
	State          string
	RestrictToPr   int
	UpdatedSeconds int
}

type Port interface {
	GetPRs(owner, repo string, opts GetPRsOptions) ([]PullRequest, error)
	ListIssueEvents(owner, repo string, issueNumber int) ([]IssueEvent, error)
}
