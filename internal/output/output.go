package output

import (
	"encoding/json"
	"fmt"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/gh"
)

type PROutput struct {
	gh.PullRequest
	LabelsStats []gh.LabelEventStats `json:"labels_stats"`
}

type Output struct {
	FormatVersion string     `json:"format_version"`
	PRs           []PROutput `json:"prs"`
}

func New() *Output {
	return &Output{
		FormatVersion: "1.0",
		PRs:           []PROutput{},
	}
}

func (output *Output) AddPROutput(pr *gh.PullRequest, stats []gh.LabelEventStats) {
	output.PRs = append(output.PRs, PROutput{
		PullRequest: *pr,
		LabelsStats: stats,
	})
}

func (o *Output) Print() error {
	output, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
