package output

import (
	"encoding/json"
	"fmt"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/repo"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/domain/stats"
)

type PROutput struct {
	repo.PullRequest
	LabelsStats []stats.LabelEventStats `json:"labels_stats"`
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

func (output *Output) AddPROutput(pr repo.PullRequest, stats []stats.LabelEventStats) {
	output.PRs = append(output.PRs, PROutput{
		PullRequest: pr,
		LabelsStats: stats,
	})
}

func (o *Output) Marshall() ([]byte, error) {
	return json.MarshalIndent(o, "", "    ")
}

func (o *Output) Print() error {
	output, err := o.Marshall()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
