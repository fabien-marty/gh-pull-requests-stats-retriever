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

type SelectPR struct {
	State          string `toml:"state" validate:"eq=open,closed,all"`
	UpdatedSeconds int    `toml:"updated_seconds"`
}

type BasicStats struct {
	Enabled bool `toml:"enabled"`
}

type LabelsStats struct {
	Enabled bool       `toml:"enabled"`
	Labels  [][]string `toml:"labels"`
}

type Config struct {
	Path        string
	Owner       string `toml:"owner" validate:"required"`
	Repo        string `toml:"repo" validate:"required"`
	Token       string
	SelectPRs   []SelectPR  `toml:"select_prs"`
	BasicStats  BasicStats  `toml:"basic_stats"`
	LabelsStats LabelsStats `toml:"labels_stats"`
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

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}
