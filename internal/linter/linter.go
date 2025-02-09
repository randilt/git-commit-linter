package linter

import (
	"fmt"
	"regexp"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/git"
)

type Linter struct {
    config *config.Config
}

func New(cfg *config.Config) *Linter {
    return &Linter{config: cfg}
}

func (l *Linter) LintCommits(commitRange string) error {
    commits, err := git.GetCommits(commitRange)
    if err != nil {
        return fmt.Errorf("failed to get commits: %w", err)
    }

    var hasErrors bool
    for _, commit := range commits {
        if err := l.lintCommit(commit); err != nil {
            fmt.Printf("Commit %s: %v\n", commit.Hash[:8], err)
            hasErrors = true
        }
    }

    if hasErrors {
        return fmt.Errorf("some commits failed linting")
    }
    return nil
}

func (l *Linter) lintCommit(commit git.Commit) error {
    // Regular expression for commit message format
    pattern := `^([\w]+)(?:\(([\w-]+)\))?: (.+)$`
    re := regexp.MustCompile(pattern)

    matches := re.FindStringSubmatch(commit.Message)
    if matches == nil {
        return fmt.Errorf("invalid commit message format. Expected: type(scope): message")
    }

    commitType := matches[1]
    scope := matches[2]
    message := matches[3]

    // Check commit type
    validType := false
    for _, t := range l.config.Types {
        if commitType == t {
            validType = true
            break
        }
    }
    if !validType {
        return fmt.Errorf("invalid commit type '%s'. Allowed types: %v", commitType, l.config.Types)
    }

    // Check scope if required
    if l.config.Rules.RequireScope && scope == "" {
        return fmt.Errorf("scope is required")
    }

    // Check message length
    if len(message) > l.config.Rules.MaxMessageLength {
        return fmt.Errorf("message exceeds maximum length of %d characters", l.config.Rules.MaxMessageLength)
    }

    return nil
}
