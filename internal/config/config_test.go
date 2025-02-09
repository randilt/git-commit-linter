package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
    // Test default config
    cfg, err := Load("")
    if err != nil {
        t.Errorf("Load() error = %v", err)
        return
    }
    if len(cfg.Types) == 0 {
        t.Error("Default config should have predefined types")
    }

    // Test custom config
    testConfig := `
types:
  - feat
  - fix
rules:
  require_scope: true
  max_message_length: 50
`
    tmpfile, err := os.CreateTemp("", "config*.yaml")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.Write([]byte(testConfig)); err != nil {
        t.Fatal(err)
    }
    if err := tmpfile.Close(); err != nil {
        t.Fatal(err)
    }

    cfg, err = Load(tmpfile.Name())
    if err != nil {
        t.Errorf("Load() error = %v", err)
        return
    }
    if len(cfg.Types) != 2 {
        t.Errorf("Expected 2 types, got %d", len(cfg.Types))
    }
    if !cfg.Rules.RequireScope {
        t.Error("RequireScope should be true")
    }
    if cfg.Rules.MaxMessageLength != 50 {
        t.Errorf("Expected MaxMessageLength = 50, got %d", cfg.Rules.MaxMessageLength)
    }
}