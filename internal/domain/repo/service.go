package repo

type Service struct {
	adapter Port
	owner   string
	repo    string
}

func NewService(adapter Port, owner string, repo string) *Service {
	return &Service{
		adapter: adapter,
		owner:   owner,
		repo:    repo,
	}
}

func (s *Service) GetPRs(optss []GetPRsOptions) ([]PullRequest, error) {
	alreadySeen := map[int]bool{}
	res := []PullRequest{}
	for _, opts := range optss {
		prs, err := s.adapter.GetPRs(s.owner, s.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			if _, ok := alreadySeen[pr.Number]; ok {
				continue
			}
			alreadySeen[pr.Number] = true
			res = append(res, pr)
		}
	}
	return res, nil
}

func (s *Service) ListIssueTimeline(issueNumber int) ([]IssueTimeline, error) {
	return s.adapter.ListIssueTimeline(s.owner, s.repo, issueNumber)
}
