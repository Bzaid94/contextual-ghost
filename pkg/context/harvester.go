package context

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Context holds the gathered information
type Context struct {
	GitDiff   string
	GitLog    string
	EnvVars   string
	OS        string
	Arch      string
}

// Harvester collects the context
type Harvester struct{}

func NewHarvester() *Harvester {
	return &Harvester{}
}

// Collect gathers all context information
func (h *Harvester) Collect() Context {
	return Context{
		GitDiff: h.getGitDiff(),
		GitLog:  h.getGitLog(),
		EnvVars: h.getEnvVars(),
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
	}
}

func (h *Harvester) getGitDiff() string {
	cmd := exec.Command("git", "diff", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return "git diff failed or not a git repo"
	}
	return strings.TrimSpace(string(out))
}

func (h *Harvester) getGitLog() string {
	cmd := exec.Command("git", "log", "-n", "3", "--oneline")
	out, err := cmd.Output()
	if err != nil {
		return "git log failed"
	}
	return strings.TrimSpace(string(out))
}

func (h *Harvester) getEnvVars() string {
	// Filter for interesting env vars to avoid noise/secrets
	interesting := []string{
		"NODE_ENV", "GOOS", "GOARCH", "SHELL", "TERM", "PATH",
		"PYTHONPATH", "GOPATH", "CARGO_HOME", "npm_config_user_agent",
	}
	
	var sb strings.Builder
	for _, key := range interesting {
		if val := os.Getenv(key); val != "" {
			sb.WriteString(fmt.Sprintf("%s=%s\n", key, val))
		}
	}
	
	// Also add any generally "relevant" ones found in os.Environ() if needed,
	// but an allowlist is safer for a prototype.
	
	return sb.String()
}
