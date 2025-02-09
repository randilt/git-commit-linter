package linter

import (
	"fmt"
	"regexp"
	"strings"

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
            fmt.Printf("\nCommit %s: %v\n", commit.Hash[:8], err)
            // Print fix instructions
            l.printFixInstructions(commit, err)
            hasErrors = true
        }
    }

    if hasErrors {
        return fmt.Errorf("\nsome commits failed linting - see fix instructions above")
    }
    return nil
}

func (l *Linter) printFixInstructions(commit git.Commit, lintErr error) {
    fmt.Println("\nTo fix this issue:")

    // Check if this is the latest commit
    isLatestCommit := strings.Contains(commit.Hash, "HEAD")
    
    if isLatestCommit {
        fmt.Println("This is your latest commit. You can fix it using:")
        fmt.Printf("  git commit --amend -m \"your-new-message\"\n")
    } else {
        fmt.Printf("This is an older commit. You can fix it using interactive rebase\nuse: git rebase -i %s~1\n", commit.Hash[:8])
    }

    fmt.Println("\nExample of valid commit message format:")
    fmt.Printf("  feat(scope): your message (max %d chars)\n", l.config.Rules.MaxMessageLength)
    
    // Print allowed types
    fmt.Println("\nAllowed types:", strings.Join(l.config.Types, ", "))
}

func (l *Linter) lintCommit(commit git.Commit) error {
    // Regular expression for commit message format
    pattern := `^([\w]+)(?:\(([\w-]+)\))?: (.+)$`
    re := regexp.MustCompile(pattern)

    matches := re.FindStringSubmatch(commit.Message)
    if matches == nil {
        return fmt.Errorf("invalid commit message format - must follow pattern: type(scope): message")
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
        return fmt.Errorf("invalid commit type '%s'", commitType)
    }

    // Check scope if required
    if l.config.Rules.RequireScope && scope == "" {
        return fmt.Errorf("scope is required")
    }

    // Check message length
    if len(message) > l.config.Rules.MaxMessageLength {
        return fmt.Errorf("message length is %d chars (maximum is %d)", 
            len(message), l.config.Rules.MaxMessageLength)
    }

    return nil
}