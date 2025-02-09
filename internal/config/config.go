package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Types  []string `yaml:"types"`
    Scopes []string `yaml:"scopes,omitempty"`
    Rules  struct {
        RequireScope    bool `yaml:"require_scope"`
        MaxMessageLength int  `yaml:"max_message_length"`
    } `yaml:"rules"`
}

func Load(path string) (*Config, error) {
    if path == "" {
        // Load default config
        return &Config{
            Types: []string{"feat", "fix", "docs", "style", "refactor", "test", "chore"},
            Rules: struct {
                RequireScope    bool `yaml:"require_scope"`
                MaxMessageLength int  `yaml:"max_message_length"`
            }{
                RequireScope:    false,
                MaxMessageLength: 72,
            },
        }, nil
    }

    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}