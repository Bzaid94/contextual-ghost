package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bdiaz/contextual-ghost/pkg/context"
	"github.com/bdiaz/contextual-ghost/pkg/runner"
	"github.com/bdiaz/contextual-ghost/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ghost <command> [args...]")
		os.Exit(0)
	}

	commandArgs := os.Args[1:]

	// 1. Run the command
	r := runner.NewRunner()
	exitCode, stderr, err := r.Run(commandArgs)

	if err != nil {
		if exitCode == -1 {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(127)
		}
	}

	if exitCode == 0 {
		os.Exit(0)
	}

	// 2. Command failed! Awakening Ghost.
	if err := checkDependencies(); err != nil {
		fmt.Printf("\n%s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("👻 Ghost cannot help you yet!"))
		fmt.Printf("%v\n", err)
		os.Exit(exitCode)
	}

	// 3. Gather Context
	harvester := context.NewHarvester()
	ctx := harvester.Collect()

	// 4. Start UI
	fullCommand := strings.Join(commandArgs, " ")
	model := ui.NewModel(ctx, stderr, fullCommand)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error within the Ghost: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func checkDependencies() error {
	// Check gh
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed.\nPlease install it: https://cli.github.com/")
	}

	// Check if 'gh copilot' command is available (either as extension or built-in)
	cmd := exec.Command("gh", "copilot", "--help")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("GitHub Copilot is not available in your gh CLI.\nPlease ensure you have access to Copilot and the CLI is authenticated: Run 'gh auth login'.")
	}

	return nil
}
