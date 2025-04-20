package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/profile"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/workflow"

	// "github.com/charmbracelet/bubbles/key"
	"github.com/jjournet/tgr/ghuser"
)

var (
	P          *tea.Program
	User       *ghuser.GHUser
	Pr         *profile.Profile
	Repo       *repository.Repository
	Workflow   *workflow.Workflow
	WindowSize tea.WindowSizeMsg
	Path       []repository.RepoElement
)

var (
	DocStyle = lipgloss.NewStyle().Margin(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	TopBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	// subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	MainStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	CommandStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	BaseTableStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Align(lipgloss.Left).Padding(0, 1).
			Foreground(lipgloss.Color("#77c2f9"))

	HighlightedLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#77c2f9")).Padding(0, 1)
)
