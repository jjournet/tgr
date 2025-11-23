package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WindowSize stores the current terminal window size
// This is the only piece of "global" state we keep for UI purposes
var WindowSize tea.WindowSizeMsg

// Styling constants
var (
	DocStyle = lipgloss.NewStyle().Margin(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	TopBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	MainStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	CommandStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	BaseTableStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Align(lipgloss.Left).Padding(0, 1).
			Foreground(lipgloss.Color("#77c2f9"))

	HighlightedLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#77c2f9")).Padding(0, 1)

	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	NoStyle      = lipgloss.NewStyle()
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)
