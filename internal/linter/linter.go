package linter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/git"
)

type Linter struct {
    config *config.Config
}

type LintError struct {
    CommitHash string
    Message    string
    FixSteps   string
}

func New(cfg *config.Config) *Linter {
    return &Linter{config: cfg}
}

// LintCommits lints the commit messages in the given range
// and returns an error if any commit fails the linting rules.
// to use in cli, run `git-commit-linter --config=config.yaml --check="HEAD~5..HEAD"` (checks last 5 commits)
func (l *Linter) LintCommits(commitRange string) error {
    commits, err := git.GetCommits(commitRange)
    if err != nil {
        return fmt.Errorf("failed to get commits: %w", err)
    }

    var lintErrors []LintError
    for _, commit := range commits {
        if err := l.lintCommit(commit); err != nil {
            lintError := LintError{
                CommitHash: commit.Hash[:8],
                Message:    err.Error(),
                FixSteps:   l.getFixInstructions(commit),
            }
            lintErrors = append(lintErrors, lintError)
        }
    }

    if len(lintErrors) > 0 {
        var output bytes.Buffer

        // Print summary header
        output.WriteString("\nLinting Issues Found:\n")
        output.WriteString("==================\n\n")

        // Print each error with its fix instructions
        for _, err := range lintErrors {
            output.WriteString(fmt.Sprintf("Commit %s: %s\n", err.CommitHash, err.Message))
            output.WriteString(err.FixSteps)
            output.WriteString("\n")
        }

        // Print reference information once at the end
        output.WriteString("\nReference Information:\n")
        output.WriteString("====================\n")
        output.WriteString(fmt.Sprintf("Valid commit format: type(scope): message (max %d chars)\n", l.config.Rules.MaxMessageLength))
        output.WriteString(fmt.Sprintf("Allowed types: %s\n", strings.Join(l.config.Types, ", ")))

        fmt.Print(output.String())
        return fmt.Errorf("commits failed linting - please fix the issues above")
    }

    fmt.Println("All commits passed linting!")
    return nil
}

func (l *Linter) getFixInstructions(commit git.Commit) string {
    var instructions strings.Builder
    instructions.WriteString("Fix Instructions:\n")

	// Check if this is the latest commit
    isLatestCommit := strings.Contains(commit.Hash, "HEAD")
    if isLatestCommit {
        instructions.WriteString("- Latest commit: Use amend\n")
        instructions.WriteString("  git commit --amend -m \"type(scope): your message\"\n")
    } else {
        instructions.WriteString("- Older commit: Use interactive rebase\n")
        instructions.WriteString(fmt.Sprintf("  git rebase -i %s~1\n", commit.Hash[:8]))
        instructions.WriteString("  Change 'pick' to 'reword' for the target commit\n")
    }

    return instructions.String()
}

func (l *Linter) lintCommit(commit git.Commit) error {
    pattern := `^([\w]+)(?:\(([\w-]+)\))?: (.+)$`
    re := regexp.MustCompile(pattern)

    matches := re.FindStringSubmatch(commit.Message)
    if matches == nil {
        return fmt.Errorf("invalid format")
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
        return fmt.Errorf("invalid type '%s'", commitType)
    }

    // Check scope if required
    if l.config.Rules.RequireScope && scope == "" {
        return fmt.Errorf("scope is required")
    }

    // Check message length
    if len(message) > l.config.Rules.MaxMessageLength {
        return fmt.Errorf("message too long (%d chars, max %d)", 
            len(message), l.config.Rules.MaxMessageLength)
    }

    return nil
}