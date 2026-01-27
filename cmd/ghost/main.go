package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bdiaz/ghost/pkg/context"
	"github.com/bdiaz/ghost/pkg/runner"
	"github.com/bdiaz/ghost/pkg/ui"
	"github.com/charmbracelet/bubbletea"
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

	// If command succeeded (exitCode 0) or failed to start (err != nil but not exit error), just exit
	// If err is not nil, it might be that the command was not found.
	if err != nil {
		// If the command itself wasn't found or couldn't start, we might print an error and exit.
		// Or should asking Ghost about it?
		// Typically if exec fails to start, it returns strict error.
		// If it started and failed, it returns ExitError (handled inside Run which returns exitCode).
		
		// If exitCode is -1, it means execution failed to start.
		if exitCode == -1 {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(127)
		}
	}
	
	if exitCode == 0 {
		os.Exit(0)
	}

	// 2. Command failed! Awakening Ghost.
	// But first, check dependencies to see if we CAN help.
	if err := checkDependencies(); err != nil {
		fmt.Printf("\n%s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Render("👻 Ghost cannot help you yet!"))
		fmt.Printf("%v\n", err)
		os.Exit(exitCode)
	}
	
	// 3. Gather Context
	// We run this in main goroutine before UI or inside UI?
	// The prompt was "The Context Harvester... Ante un fallo, recolectar en paralelo".
	// We'll collect it here quickly.
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
	
	// Check copilot extension
	cmd := exec.Command("gh", "extension", "list")
	out, err := cmd.Output()
	if err != nil {
		// If this fails, maybe not authenticated?
		return fmt.Errorf("Could not check gh extensions. Are you authenticated? Run 'gh auth login'.")
	}
	
	if !strings.Contains(string(out), "copilot") {
		return fmt.Errorf("GitHub Copilot extension is not installed.\nPlease run: gh extension install github/gh-copilot")
	}
	
	return nil
}
