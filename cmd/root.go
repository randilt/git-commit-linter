package cmd

import (
	"fmt"
	"os"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/git"
	"github.com/randilt/git-commit-linter/internal/linter"
	"github.com/spf13/cobra"
)

var (
	version     string
	commit      string
	date        string
	configPath  string
	commitRange string

	rootCmd = &cobra.Command{
		Use:   "git-commit-linter",
		Short: "A tool to lint Git commit messages",
		Long: `Git Commit Linter ensures your commit messages follow standardized formats.
Example: git-commit-linter --config=config.yaml --check="HEAD~5..HEAD"`,
		RunE: runLinter,
	}

	installHookCmd = &cobra.Command{
		Use:   "install-hook",
		Short: "Install git commit-msg hook",
		Long: `Installs a git commit-msg hook that will automatically lint 
commit messages before they are committed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return git.InstallHook()
		},
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
	rootCmd.AddCommand(installHookCmd)
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
		return fmt.Errorf("failed to load config: %w", err)
	}

	l := linter.New(cfg)
	if err := l.LintCommitMessageFile(args[0]); err != nil {
		if _, ok := err.(*linter.ValidationError); ok {
			// For validation errors, just exit with status code 1
			os.Exit(1)
		}
		// For other errors, return them to cobra
		return err
	}
	return nil
}
