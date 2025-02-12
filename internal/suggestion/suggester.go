package suggestion

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type KeywordType struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Keywords    []string `yaml:"keywords"`
}

type KeywordsConfig struct {
	CommitTypes  []KeywordType `yaml:"commit_types"`
	CommitScopes []KeywordType `yaml:"commit_scopes"`
}

func LoadKeywords() (*KeywordsConfig, error) {
	// Try to find common_keywords.yaml in various locations
	locations := []string{
		"config/common_keywords.yaml",
		"../config/common_keywords.yaml",
		"../../config/common_keywords.yaml",
	}

	var configData []byte
	var err error

	for _, loc := range locations {
		configData, err = os.ReadFile(loc)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var config KeywordsConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SuggestCommitMessage(analysis *ChangeAnalysis) (*CommitSuggestion, error) {
	if analysis == nil || len(analysis.FilesChanged) == 0 {
		return nil, nil
	}

	keywords, err := LoadKeywords()
	if err != nil {
		return nil, err
	}

	suggestions := make(map[string]float64)

	// Score each commit type based on the analysis
	for _, commitType := range keywords.CommitTypes {
		score := 0.0

		// Check file patterns
		if commitType.Name == "docs" && analysis.FilePatterns[".md"] > 0 {
			score += 0.8
		}
		if commitType.Name == "test" && analysis.FilePatterns["test"] > 0 {
			score += 0.8
		}
		if commitType.Name == "style" && (analysis.FilePatterns[".css"] > 0 || analysis.FilePatterns[".scss"] > 0) {
			score += 0.6
		}

		// Analyze content changes
		for _, keyword := range commitType.Keywords {
			pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)

			for _, line := range analysis.AddedLines {
				if pattern.MatchString(line) {
					score += 0.3
				}
			}
			for _, line := range analysis.RemovedLines {
				if pattern.MatchString(line) {
					score += 0.2
				}
			}
		}

		if score > 0 {
			suggestions[commitType.Name] = score
		}
	}

	// Find the best match
	var bestType string
	var bestScore float64
	for typ, score := range suggestions {
		if score > bestScore {
			bestType = typ
			bestScore = score
		}
	}

	if bestScore < 0.3 {
		// Not confident enough to make a suggestion
		return nil, nil
	}

	// Generate a message based on the changes
	message := generateCommitMessage(analysis, bestType)

	return &CommitSuggestion{
		Type:    bestType,
		Message: message,
		Score:   bestScore,
	}, nil
}

func generateCommitMessage(analysis *ChangeAnalysis, commitType string) string {
	if len(analysis.FilesChanged) == 0 {
		return ""
	}

	// Get the most common file type
	mostCommonExt := ""
	maxCount := 0
	for ext, count := range analysis.FilePatterns {
		if count > maxCount {
			mostCommonExt = ext
			maxCount = count
		}
	}

	switch commitType {
	case "feat":
		return "add " + describeMajorChange(analysis)
	case "fix":
		return "fix " + describeMajorChange(analysis)
	case "docs":
		if len(analysis.FilesChanged) == 1 {
			return "update documentation for " + analysis.FilesChanged[0]
		}
		return "update documentation"
	case "test":
		return "add tests for " + describeMajorChange(analysis)
	case "style":
		return "improve " + mostCommonExt + " styling"
	case "refactor":
		return "refactor " + describeMajorChange(analysis)
	default:
		return describeMajorChange(analysis)
	}
}

func describeMajorChange(analysis *ChangeAnalysis) string {
	if len(analysis.FilesChanged) == 1 {
		return strings.TrimSuffix(filepath.Base(analysis.FilesChanged[0]), filepath.Ext(analysis.FilesChanged[0]))
	}

	dir := filepath.Dir(analysis.FilesChanged[0])
	if dir == "." {
		return "multiple files"
	}
	return filepath.Base(dir) + " module"
}
