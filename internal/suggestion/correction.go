package suggestion

import (
	"fmt"
	"sort"
	"strings"
)

type CommitCorrection struct {
    Type    string
    Scope   string
    Message string
    Score   float64
}
type KeywordType struct {
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Keywords    []string `yaml:"keywords"`
}

type KeywordsConfig struct {
    CommitTypes  []KeywordType `yaml:"commit_types"`
    CommitScopes []KeywordType `yaml:"commit_scopes"`
}

// SuggestCorrection analyzes an invalid commit message and suggests corrections
func SuggestCorrection(message string, config *KeywordsConfig) (*CommitCorrection, error) {
    // Clean the message
    message = strings.TrimSpace(message)
    
    // Skip if message is empty
    if message == "" {
        return nil, fmt.Errorf("empty commit message")
    }

    // Calculate scores for each type based on keywords
    typeScores := make(map[string]float64)
    for _, commitType := range config.CommitTypes {
        score := calculateTypeScore(message, commitType.Keywords)
        if score > 0 {
            typeScores[commitType.Name] = score
        }
    }

    // Find the best matching type
    bestType := findBestMatch(typeScores)
    if bestType == "" {
        bestType = "chore" // Default to chore if no clear match
    }

    // Calculate scores for each scope based on keywords
    scopeScores := make(map[string]float64)
    for _, scope := range config.CommitScopes {
        score := calculateScopeScore(message, scope.Keywords)
        if score > 0 {
            scopeScores[scope.Name] = score
        }
    }

    // Find the best matching scope
    bestScope := findBestMatch(scopeScores)

    return &CommitCorrection{
        Type:    bestType,
        Scope:   bestScope,
        Message: message,
        Score:   typeScores[bestType],
    }, nil
}

func calculateTypeScore(message string, keywords []string) float64 {
    messageLower := strings.ToLower(message)
    var score float64

    // Check if message starts with any keyword (higher weight)
    for _, keyword := range keywords {
        if strings.HasPrefix(messageLower, strings.ToLower(keyword)) {
            score += 0.6
        }
    }

    // Check for keyword presence anywhere in the message
    for _, keyword := range keywords {
        if strings.Contains(messageLower, strings.ToLower(keyword)) {
            score += 0.3
        }
    }

    return score
}

func calculateScopeScore(message string, keywords []string) float64 {
    messageLower := strings.ToLower(message)
    var score float64

    for _, keyword := range keywords {
        if strings.Contains(messageLower, strings.ToLower(keyword)) {
            score += 0.4
        }
    }

    return score
}

func findBestMatch(scores map[string]float64) string {
    if len(scores) == 0 {
        return ""
    }

    type scoreEntry struct {
        key   string
        score float64
    }

    var entries []scoreEntry
    for k, v := range scores {
        entries = append(entries, scoreEntry{k, v})
    }

    sort.Slice(entries, func(i, j int) bool {
        return entries[i].score > entries[j].score
    })

    if entries[0].score >= 0.3 { // Minimum confidence threshold
        return entries[0].key
    }

    return ""
}