package main

import (
	"fmt"
	"os"

	"github.com/randilt/git-commit-linter/cmd"
)

// Version information (filled by GoReleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
