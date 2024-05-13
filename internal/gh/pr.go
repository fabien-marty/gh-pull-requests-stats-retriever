package gh

import (
	"context"
	"log/slog"
	"time"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/config"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/log"
	"github.com/google/go-github/v62/github"
)

type PullRequest struct {
	Number    int              `json:"number"`
	CreatedAt *time.Time       `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at"`
	ClosedAt  *time.Time       `json:"closed_at"`
	MergedAt  *time.Time       `json:"merged_at"`
	State     PullRequestState `json:"state"`
	Labels    []string         `json:"current_labels"`
}

func NewPullRequest(pr *github.PullRequest) *PullRequest {
	state, err := ParsePullRequestState(pr.GetState())
	if err != nil {
		panic(err)
	}
	labels := []string{}
	for _, lbl := range pr.Labels {
		name := lbl.GetName()
		if name == "" {
			continue
		}
		labels = append(labels, name)
	}
	return &PullRequest{
		Number:    pr.GetNumber(),
		State:     state,
		CreatedAt: pr.CreatedAt.GetTime(),
		UpdatedAt: pr.UpdatedAt.GetTime(),
		ClosedAt:  pr.ClosedAt.GetTime(),
		MergedAt:  pr.MergedAt.GetTime(),
		Labels:    labels,
	}
}

func GetPRs(client *Client, config *config.Config) ([]*PullRequest, error) {
	logger := log.GetLogger().With(slog.String("owner", config.Owner), slog.String("repo", config.Repo))
	options := github.PullRequestListOptions{
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			Page: 1,
		},
	}
	alreadySeen := map[int64]bool{}
	res := []*PullRequest{}
	for _, selectPR := range config.SelectPRs {
		options.State = selectPR.State
	out:
		for {
			logger.Info("Fetching PRs...", slog.String("state", selectPR.State), slog.Int("page", options.Page))
			prs, resp, err := client.client.PullRequests.List(context.Background(), config.Owner, config.Repo, &options)
			if err != nil {
				return nil, err
			}
			for _, pr := range prs {
				if config.RestrictToPr > 0 && pr.GetNumber() != config.RestrictToPr {
					continue
				}
				if _, ok := alreadySeen[pr.GetID()]; ok {
					continue
				}
				if pr.Number == nil {
					logger.Warn("PR has no number", slog.Int64("id", pr.GetID()))
					continue
				}
				since := int(time.Since(*pr.UpdatedAt.GetTime()).Seconds())
				if selectPR.UpdatedSeconds > 0 && since > selectPR.UpdatedSeconds {
					logger.Debug("PR too old", slog.Int("number", pr.GetNumber()), slog.Int("since", since), slog.Int("threshold", selectPR.UpdatedSeconds))
					break
				}
				alreadySeen[pr.GetID()] = true
				res = append(res, NewPullRequest(pr))
				logger.Info("PR fetched", slog.Int("number", pr.GetNumber()))
			}
			if resp.NextPage == 0 {
				break out
			}
			options.Page = resp.NextPage
		}
	}
	logger.Info("PRs fetched", slog.Int("count", len(res)))
	return res, nil
}
