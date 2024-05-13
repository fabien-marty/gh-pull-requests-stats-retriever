package gh

import (
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/config"
	"github.com/google/go-github/v62/github"
)

type Client struct {
	client *github.Client
}

func GetClient(config *config.Config) *Client {
	client := github.NewClient(nil)
	if config.Token != "" {
		client = client.WithAuthToken(config.Token)
	}
	return &Client{client: client}
}
