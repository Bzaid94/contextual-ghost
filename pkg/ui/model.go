package ui

import (
	"fmt"
	"strings"

	"time"

	"github.com/bdiaz/contextual-ghost/pkg/bridge"
	"github.com/bdiaz/contextual-ghost/pkg/context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	headerStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2)

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type errMsg error

type explanationMsg string

type Model struct {
	context  context.Context
	errorLog string
	command  string
	bridge   *bridge.Bridge

	explanation string
	loading     bool
	err         error

	// Spinner (simplified for now, usually requires a spinner bubble)
	spinnerFrame int
}

func NewModel(ctx context.Context, errorLog, command string) Model {
	return Model{
		context:  ctx,
		errorLog: errorLog,
		command:  command,
		bridge:   bridge.NewBridge(),
		loading:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchExplanation,
		m.tickSpinner,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case explanationMsg:
		m.loading = false
		m.explanation = string(msg)
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg
		return m, nil

	case spinnerTickMsg:
		if !m.loading {
			return m, nil
		}
		m.spinnerFrame++
		return m, m.tickSpinner
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.loading {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := frames[m.spinnerFrame%len(frames)]
		return fmt.Sprintf("\n %s Consulting the Ghost...\n", spinnerStyle.Render(frame))
	}

	// Structured view

	doc := strings.Builder{}

	doc.WriteString(headerStyle.Render("👻 Ghost Agent Analysis"))
	doc.WriteString("\n\n")

	// Explanation box
	doc.WriteString(boxStyle.Render(m.explanation))
	doc.WriteString("\n\n")

	doc.WriteString(lipgloss.NewStyle().Foreground(subtle).Render("Press 'q' to dismiss"))

	return doc.String()
}

func (m Model) fetchExplanation() tea.Msg {
	res, err := m.bridge.Ask(m.context, m.errorLog, m.command)
	if err != nil {
		return errMsg(err)
	}
	return explanationMsg(res)
}

type spinnerTickMsg struct{}

func (m Model) tickSpinner() tea.Msg {
	return tea.Tick(time.Second/10, func(_ time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}
