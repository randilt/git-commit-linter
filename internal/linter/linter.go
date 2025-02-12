package linter

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/git"
	"github.com/randilt/git-commit-linter/internal/suggestion"
	"github.com/randilt/git-commit-linter/internal/ui"
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

// LintCommitMessage lints a single commit message from a string
func (l *Linter) LintCommitMessage(message string) error {
	// Create a temporary commit object to reuse existing logic
	tempCommit := git.Commit{
		Hash:    "UNCOMMITTED",
		Message: message,
	}

	if err := l.lintCommit(tempCommit); err != nil {
		ui.Section("Linting Issues Found")
		ui.Error(err.Error())

		ui.Section("Reference Information")
		ui.Info(fmt.Sprintf("Valid commit format: %s",
			ui.Bold(fmt.Sprintf("type(scope): message (max %d chars)", l.config.Rules.MaxMessageLength))))
		ui.Info(fmt.Sprintf("Allowed types: %s",
			ui.Bold(strings.Join(l.config.Types, ", "))))

		return fmt.Errorf("commit message failed linting - please fix the issues above")
	}

	ui.Success("Commit message passed linting!")
	return nil
}

// LintCommitMessageFile lints a commit message from a file path
func (l *Linter) LintCommitMessageFile(filepath string) error {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read commit message file: %w", err)
	}

	// Clean the message - remove comments and empty lines
	lines := strings.Split(string(content), "\n")
	var messageLines []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") && len(strings.TrimSpace(line)) > 0 {
			messageLines = append(messageLines, line)
		}
	}
	message := strings.Join(messageLines, "\n")

	return l.LintCommitMessage(strings.TrimSpace(message))
}

func (l *Linter) SuggestMessageCorrection(message string) (string, error) {
	// Load keywords configuration
	keywords, err := suggestion.LoadKeywords()
	if err != nil {
		return "", fmt.Errorf("failed to load keywords: %w", err)
	}

	// Get suggestion
	correction, err := suggestion.SuggestCorrection(message, keywords)
	if err != nil {
		return "", err
	}

	// Format the suggested commit message
	var suggestionBuilder strings.Builder

	// Add type
	suggestionBuilder.WriteString(correction.Type)

	// Add scope if present and allowed
	if correction.Scope != "" && !l.config.Rules.RequireScope {
		suggestionBuilder.WriteString(fmt.Sprintf("(%s)", correction.Scope))
	}

	// Add message
	suggestionBuilder.WriteString(": ")
	suggestionBuilder.WriteString(correction.Message)

	return suggestionBuilder.String(), nil
}

// LintCommits lints the commit messages in the given range and returns an error if any commit fails the linting rules.
//
// The function accepts a Git commit range (e.g., "HEAD~5..HEAD") and validates each commit message against
// the configured linting rules. It performs the following checks:
//   - Commit message format validation
//   - Commit type verification
//   - Scope requirement check (if enabled)
//   - Message length validation
//
// For any commits that fail validation, it generates detailed error messages including:
//   - The commit hash
//   - The specific validation failure
//   - Step-by-step instructions for fixing the commit
//
// Example usage:
//
//	git-commit-linter --config=config.yaml --check="HEAD~5..HEAD"
//
// Returns nil if all commits pass validation, or error details if any commits fail.
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
		ui.Section("Linting Issues Found")

		// Print each error with its fix instructions
		for _, err := range lintErrors {
			ui.Error(fmt.Sprintf("Commit %s: %s", ui.Bold(err.CommitHash), err.Message))
			ui.CodeBlock(err.FixSteps)
		}

		ui.Section("Reference Information")
		ui.Info(fmt.Sprintf("Valid commit format: %s",
			ui.Bold(fmt.Sprintf("type(scope): message (max %d chars)", l.config.Rules.MaxMessageLength))))
		ui.Info(fmt.Sprintf("Allowed types: %s",
			ui.Bold(strings.Join(l.config.Types, ", "))))

		ui.Error("Some commits failed linting - please fix the issues above")
		return nil
	}

	ui.Success("All commits passed linting!")
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
		// Message format is invalid, try to suggest a correction
		suggestion, err := l.SuggestMessageCorrection(commit.Message)
		if err == nil && suggestion != "" {
			return fmt.Errorf("invalid format. Did you mean: %s", ui.Bold(suggestion))
		}
		return fmt.Errorf("invalid format")
	}

	// Rest of the existing validation logic remains the same...
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
		// If type is invalid, try to suggest a correction
		suggestion, err := l.SuggestMessageCorrection(commit.Message)
		if err == nil && suggestion != "" {
			return fmt.Errorf("invalid type '%s'. Did you mean: %s", commitType, ui.Bold(suggestion))
		}
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
