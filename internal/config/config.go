package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"

	"github.com/fabien-marty/gh-pull-requests-stats-retriever/internal/log"
	"github.com/pelletier/go-toml/v2"
)

// SelectPR represents the configuration of a "PR to select" block
type SelectPR struct {
	State          string `toml:"state" validate:"eq=open,closed,all"`
	UpdatedSeconds int    `toml:"updated_seconds"`
}

// LabelsStats represents the configuration of the labels stats
type LabelsStats struct {
	Enabled bool       `toml:"enabled"`
	Labels  [][]string `toml:"labels"`
}

// Config represents the configuration of the application
type Config struct {
	Path         string
	Owner        string `toml:"owner" validate:"required"`
	Repo         string `toml:"repo" validate:"required"`
	RestrictToPr int
	Token        string
	SelectPRs    []SelectPR  `toml:"select_prs"`
	LabelsStats  LabelsStats `toml:"labels_stats"`
}

func readConfig(path string) ([]byte, error) {
	log.GetLogger().Debug("Reading config file...", slog.String("path", path))
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("the given config path: %s does not exist: %w", path, err)
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("the given config path: %s is a directory", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return content, nil
}

// Parse parses the config file at the given path
func Parse(path string) (*Config, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	var cfg Config = Config{Path: absolutePath}
	doc, err := readConfig(absolutePath)
	if err != nil {
		return nil, err
	}
	log.GetLogger().Debug("Parsing config file...", slog.String("path", absolutePath))
	err = toml.Unmarshal(doc, &cfg)
	if err != nil {
		return nil, err
	}
	log.GetLogger().Debug("Config file parsed")
	return &cfg, nil
}

// Validate validates the current config object
func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}
