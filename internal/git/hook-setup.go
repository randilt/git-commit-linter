package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/randilt/git-commit-linter/internal/ui"
)

const hookContent = `#!/bin/sh
commit_msg_file="$1"

# Run the linter with the commit message file
git-commit-linter lint-file "$commit_msg_file" || exit 1
`

// InstallHook installs the commit-msg hook in the current git repository
func InstallHook() error {
	// Check if inside a git directory
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	hookPath := filepath.Join(".git", "hooks", "commit-msg")
	samplePath := filepath.Join(".git", "hooks", "commit-msg.sample")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Read existing hook content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("failed to read existing hook: %w", err)
		}

		// Prompt user for action
		ui.Section("Existing Hook")
		ui.Warning("A commit-msg hook already exists")
		ui.Info("Current hook content:")
		ui.CodeBlock(string(content))

		response := ui.Prompt("Do you want to overwrite it? [y/N]:")
		if response != "y" && response != "Y" {
			ui.Error("Hook installation cancelled")
			return nil
		}
	}

	// Create or overwrite the hook file
	err := os.WriteFile(hookPath, []byte(hookContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	// If sample hook exists, keep a backup
	if _, err := os.Stat(samplePath); err == nil {
		backupPath := samplePath + ".backup"
		if err := os.Rename(samplePath, backupPath); err != nil {
			ui.Warning(fmt.Sprintf("Could not backup sample hook: %v", err))
		}
	}

	ui.Success("Git commit-msg hook installed successfully!")
	return nil
}