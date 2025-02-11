package cmd

import (
	"fmt"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/linter"
	"github.com/spf13/cobra"
)

var (
	version    string
	commit     string
	date       string
	configPath string
	commitRange string

	rootCmd = &cobra.Command{
		Use:   "git-commit-linter",
		Short: "A tool to lint Git commit messages",
		Long: `Git Commit Linter ensures your commit messages follow standardized formats.
Example: git-commit-linter --config=config.yaml --check="HEAD~5..HEAD"`,
		RunE: runLinter,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("git-commit-linter version %s\n", version)
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("built at: %s\n", date)
		},
	}

	lintFileCmd = &cobra.Command{
		Use:   "lint-file [file]",
		Short: "Lint a commit message from a file",
		Args:  cobra.ExactArgs(1),
		RunE:  lintFile,
	}
)

func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to config file")
	rootCmd.PersistentFlags().StringVar(&commitRange, "check", "HEAD^..HEAD", "commit range to check")
	rootCmd.AddCommand(lintFileCmd)
	rootCmd.AddCommand(versionCmd)
}

func runLinter(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(configPath)
    if err != nil {
        return err
    }

    l := linter.New(cfg)
    return l.LintCommits(commitRange)
}

func lintFile(cmd *cobra.Command, args []string) error {
    cfg, err := config.Load(configPath)
    if err != nil {
        return err
    }

    l := linter.New(cfg)
    return l.LintCommitMessageFile(args[0])
}