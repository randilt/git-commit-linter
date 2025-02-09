package git

import (
	"os/exec"
	"strings"
)

type Commit struct {
    Hash    string
    Message string
}

func GetCommits(commitRange string) ([]Commit, error) {
    cmd := exec.Command("git", "log", "--format=%H%n%B%n---", commitRange)
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    commits := []Commit{}
    parts := strings.Split(string(output), "---\n")
    
    for _, part := range parts {
        if part == "" {
            continue
        }
        
        lines := strings.SplitN(strings.TrimSpace(part), "\n", 2)
        if len(lines) < 2 {
            continue
        }
        
        commits = append(commits, Commit{
            Hash:    lines[0],
            Message: strings.TrimSpace(lines[1]),
        })
    }

    return commits, nil
}