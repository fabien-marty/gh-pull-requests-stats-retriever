package gh

import "fmt"

type PullRequestState string

const PullRequestStateClosed PullRequestState = "closed"
const PullRequestStateOpen PullRequestState = "open"

func ParsePullRequestState(state string) (PullRequestState, error) {
	switch state {
	case "closed":
		return PullRequestStateClosed, nil
	case "open":
		return PullRequestStateOpen, nil
	default:
		return "", fmt.Errorf("invalid pull request state: %s", state)
	}
}
