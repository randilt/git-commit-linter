package linter

import (
	"testing"

	"github.com/randilt/git-commit-linter/internal/config"
	"github.com/randilt/git-commit-linter/internal/git"
)

func TestLinter_LintCommit(t *testing.T) {
    cfg := &config.Config{
        Types: []string{"feat", "fix", "docs"},
        Rules: struct {
            RequireScope    bool `yaml:"require_scope"`
            MaxMessageLength int  `yaml:"max_message_length"`
        }{
            RequireScope:    false,
            MaxMessageLength: 72,
        },
    }

    linter := New(cfg)

    tests := []struct {
        name    string
        commit  git.Commit
        wantErr bool
    }{
        {
            name: "valid commit without scope",
            commit: git.Commit{
                Hash:    "abc123",
                Message: "feat: add new feature",
            },
            wantErr: false,
        },
        {
            name: "valid commit with scope",
            commit: git.Commit{
                Hash:    "def456",
                Message: "fix(auth): fix login issue",
            },
            wantErr: false,
        },
        {
            name: "invalid type",
            commit: git.Commit{
                Hash:    "ghi789",
                Message: "invalid: something",
            },
            wantErr: true,
        },
        {
            name: "message too long",
            commit: git.Commit{
                Hash:    "jkl012",
                Message: "feat: " + string(make([]byte, 100)),
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := linter.lintCommit(tt.commit)
            if (err != nil) != tt.wantErr {
                t.Errorf("LintCommit() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}