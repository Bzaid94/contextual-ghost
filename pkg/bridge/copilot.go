package bridge

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bdiaz/ghost/pkg/context"
)

type Bridge struct{}

func NewBridge() *Bridge {
	return &Bridge{}
}

func (b *Bridge) Ask(ctx context.Context, errorLog string, command string) (string, error) {
	prompt := b.constructPrompt(ctx, errorLog, command)
	
	// We use 'gh copilot explain' with the prompt.
	// Note: We might run into issues if the prompt is too long for a single argument.
	// However, for this prototype we assume it works.
	
	cmd := exec.Command("gh", "copilot", "explain", prompt)
	
	// We probably want to capture output to show it in our UI, not stream it directly yet.
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If gh fails, we return the error output
		return "", fmt.Errorf("gh copilot failed: %s (%w)", string(output), err)
	}
	
	return string(output), nil
}

func (b *Bridge) constructPrompt(ctx context.Context, errorLog string, command string) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("I am running: %s\n", command))
	sb.WriteString(fmt.Sprintf("It failed with error:\n%s\n", errorLog))
	
	if ctx.GitDiff != "" {
		sb.WriteString(fmt.Sprintf("Recent changes in files:\n%s\n", ctx.GitDiff))
	}
	
	if ctx.GitLog != "" {
		sb.WriteString(fmt.Sprintf("Last intent:\n%s\n", ctx.GitLog))
	}
	
	if ctx.EnvVars != "" {
		sb.WriteString(fmt.Sprintf("Environment:\n%s\n", ctx.EnvVars))
	}
	
	sb.WriteString("Using GitHub Copilot CLI, explain what happened and suggest the specific command to fix it if applicable.")
	
	return sb.String()
}
