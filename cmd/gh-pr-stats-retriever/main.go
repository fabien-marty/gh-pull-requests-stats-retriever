package main

import (
	"fmt"
	"os"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/config"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/gh"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/log"
	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/output"
	"github.com/urfave/cli/v2"
)

func forceConfigFromCliContext(context *cli.Context, config *config.Config) {
	owner := context.String("owner")
	if owner != "" {
		config.Owner = owner
	}
	repo := context.String("repo")
	if repo != "" {
		config.Repo = repo
	}
	token := context.String("token")
	if token != "" {
		config.Token = token
	}
	restrictToPr := context.Int("restrict-to-pr")
	if restrictToPr != 0 {
		config.RestrictToPr = restrictToPr
	}
}

func app(context *cli.Context) error {
	logLevel := context.String("log-level")
	if err := log.SetDefaultLevelFromString(logLevel); err != nil {
		return err
	}
	configPath := context.String("config")
	config, err := config.Parse(configPath)
	if err != nil {
		return err
	}
	forceConfigFromCliContext(context, config)
	err = config.Validate()
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	log.GetLogger().Debug("Configuration looks valid")
	client := gh.GetClient(config)
	prs, err := gh.GetPRs(client, config)
	if err != nil {
		return err
	}
	out := output.New()
	for _, pr := range prs {
		var stats []gh.LabelEventStats = nil
		if config.LabelsStats.Enabled {
			stats, err = gh.GetLabelEventStats(client, config, pr)
			if err != nil {
				return err
			}
		}
		out.AddPROutput(pr, stats)
	}
	out.Print()
	return nil
}

func main() {
	app := &cli.App{
		Name:   "gh-pr-stats-retriever",
		Usage:  "Get stats from GitHub PRs and dumps them to a JSON file",
		Action: app,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "owner",
				Value:   "",
				Usage:   "",
				EnvVars: []string{"GH_PR_STATS_OWNER"},
			},
			&cli.StringFlag{
				Name:    "repo",
				Value:   "",
				Usage:   "",
				EnvVars: []string{"GH_PR_STATS_REPO"},
			},
			&cli.IntFlag{
				Name:  "restrict-to-pr",
				Value: 0,
				Usage: "restrict stats to a specific PR number",
			},
			&cli.StringFlag{
				Name:    "token",
				Value:   "",
				Usage:   "GitHub (PAT) token",
				EnvVars: []string{"GH_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   "./config.toml",
				Usage:   "Path to the configuration file",
				EnvVars: []string{"GH_PR_STATS_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "INFO",
				Usage:   "Log level to use: DEBUG, INFO, WARN, ERROR",
				EnvVars: []string{"GH_PR_STATS_LOG_LEVEL"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.GetLogger().Error(err.Error())
	}
}
