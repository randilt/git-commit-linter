package git

import (
	"fmt"
	"os"
	"path/filepath"
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
		fmt.Println("commit-msg hook already exists.")
		fmt.Println("Current hook content:")
		fmt.Println("-------------------")
		fmt.Println(string(content))
		fmt.Println("-------------------")
		fmt.Print("Do you want to overwrite it? [y/N]: ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Hook installation cancelled")
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
			fmt.Printf("Warning: Could not backup sample hook: %v\n", err)
		}
	}

	fmt.Println("Git commit-msg hook installed successfully!")
	return nil
}