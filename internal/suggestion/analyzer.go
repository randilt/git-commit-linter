package suggestion

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type CommitSuggestion struct {
	Type    string
	Scope   string
	Message string
	Score   float64
}

type ChangeAnalysis struct {
	AddedLines   []string
	RemovedLines []string
	FilesChanged []string
	FilePatterns map[string]int // Patterns like *.test.js, *.md, etc.
}

func AnalyzeChanges() (*ChangeAnalysis, error) {
	// Get staged files
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, nil // No staged changes
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Get diff content
	cmd = exec.Command("git", "diff", "--cached", "--unified=0")
	diffOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	analysis := &ChangeAnalysis{
		FilesChanged: files,
		FilePatterns: make(map[string]int),
	}

	// Analyze file patterns
	for _, file := range files {
		ext := filepath.Ext(file)
		base := filepath.Base(file)
		analysis.FilePatterns[ext]++

		// Check for specific patterns
		if strings.Contains(base, ".test.") {
			analysis.FilePatterns["test"]++
		}
		if strings.Contains(base, ".spec.") {
			analysis.FilePatterns["test"]++
		}
		if ext == ".md" || ext == ".txt" || ext == ".rst" {
			analysis.FilePatterns["docs"]++
		}
	}

	// Parse diff output
	lines := strings.Split(string(diffOutput), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			analysis.AddedLines = append(analysis.AddedLines, strings.TrimPrefix(line, "+"))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			analysis.RemovedLines = append(analysis.RemovedLines, strings.TrimPrefix(line, "-"))
		}
	}

	return analysis, nil
}
