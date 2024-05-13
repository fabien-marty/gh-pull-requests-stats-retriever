package adapters

import (
	"context"
	"log/slog"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/config"
	domain_repo "github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/log"
	"github.com/google/go-github/v62/github"
)

type repoGitHub struct {
	client *github.Client
	config *config.Config
}

func NewRepoGitHub(config *config.Config) domain_repo.Port {
	client := github.NewClient(nil)
	if config.Token != "" {
		client = client.WithAuthToken(config.Token)
	}
	return &repoGitHub{
		config: config,
		client: client,
	}
}

func newPullRequest(pr *github.PullRequest) domain_repo.PullRequest {
	labels := []string{}
	for _, lbl := range pr.Labels {
		name := lbl.GetName()
		if name == "" {
			continue
		}
		labels = append(labels, name)
	}
	return domain_repo.PullRequest{
		Number:    pr.GetNumber(),
		State:     pr.GetState(),
		CreatedAt: pr.CreatedAt.GetTime(),
		UpdatedAt: pr.UpdatedAt.GetTime(),
		ClosedAt:  pr.ClosedAt.GetTime(),
		MergedAt:  pr.MergedAt.GetTime(),
		Labels:    labels,
	}
}

func (r *repoGitHub) GetPRs(owner, repo string, opts domain_repo.GetPRsOptions) ([]domain_repo.PullRequest, error) {
	logger := log.GetLogger().With(slog.String("owner", owner), slog.String("repo", repo))
	options := github.PullRequestListOptions{
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page: 1,
		},
		State: opts.State,
	}
	res := []domain_repo.PullRequest{}
out:
	for {
		logger.Info("Fetching PRs...", slog.String("state", opts.State), slog.Int("page", options.Page))
		prs, resp, err := r.client.PullRequests.List(context.Background(), owner, repo, &options)
		if err != nil {
			return nil, err
		}
		for _, pr := range prs {
			if opts.RestrictToPr > 0 && pr.GetNumber() != opts.RestrictToPr {
				continue
			}
			if pr.Number == nil {
				logger.Warn("PR has no number", slog.Int64("id", pr.GetID()))
				continue
			}
			since := int(time.Since(*pr.UpdatedAt.GetTime()).Seconds())
			if opts.UpdatedSeconds > 0 && since > opts.UpdatedSeconds {
				logger.Debug("PR too old", slog.Int("number", pr.GetNumber()), slog.Int("since", since), slog.Int("threshold", opts.UpdatedSeconds))
				break out
			}
			res = append(res, newPullRequest(pr))
		}
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}
	logger.Debug("Fetched PRs", slog.Int("count", len(res)))
	return res, nil
}

func newIssueEvent(evt *github.IssueEvent) domain_repo.IssueEvent {
	label := ""
	lbl := evt.GetLabel()
	if lbl != nil {
		label = lbl.GetName()
	}
	return domain_repo.IssueEvent{
		Type:      domain_repo.ParseIssueEventType(evt.GetEvent()),
		CreatedAt: *evt.CreatedAt.GetTime(),
		Label:     label,
	}
}

func (r *repoGitHub) ListIssueEvents(owner, repo string, issueNumber int) ([]domain_repo.IssueEvent, error) {
	logger := log.GetLogger().With(slog.String("owner", owner), slog.String("repo", repo))
	options := &github.ListOptions{
		Page: 1,
	}
	res := []domain_repo.IssueEvent{}
	for {
		logger.Info("Fetching label issue events...", slog.Int("pr.Number", issueNumber), slog.Int("page", options.Page))
		evts, resp, err := r.client.Issues.ListIssueEvents(context.Background(), owner, repo, issueNumber, options)
		if err != nil {
			return nil, err
		}
		for _, evt := range evts {
			res = append(res, newIssueEvent(evt))
		}
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}
	logger.Debug("Fetched label issue events", slog.Int("count", len(res)))
	return res, nil
}
