package repo

import "time"

type PullRequest struct {
	Number    int        `json:"number"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ClosedAt  *time.Time `json:"closed_at"`
	MergedAt  *time.Time `json:"merged_at"`
	State     string     `json:"state"`
	Labels    []string   `json:"current_labels"`
}
