# gh-pull-requests-stats-retriever
Retrieve advanced statistics about github pull-requests (pr) for a given repo (including labels stats)

## CLI usage

```console

$ gh-pr-stats-retriever --help
NAME:
   gh-pr-stats-retriever - Get stats from GitHub PRs and dumps them to a JSON file

USAGE:
   gh-pr-stats-retriever [global options] command [command options] 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value            [$GH_PR_STATS_OWNER]
   --repo value             [$GH_PR_STATS_REPO]
   --restrict-to-pr value  restrict stats to a specific PR number (default: 0)
   --token value           GitHub (PAT) token [$GH_TOKEN]
   --config value          Path to the configuration file (default: "./config.toml") [$GH_PR_STATS_CONFIG]
   --log-level value       Log level to use: DEBUG, INFO, WARN, ERROR (default: "INFO") [$GH_PR_STATS_LOG_LEVEL]
   --help, -h              show help

```
