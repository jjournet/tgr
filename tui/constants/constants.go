package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/profile"
	"github.com/jjournet/tgr/repository"

	// "github.com/charmbracelet/bubbles/key"
	"github.com/jjournet/tgr/ghuser"
)

var (
	P          *tea.Program
	User       *ghuser.GHUser
	Pr         *profile.Profile
	Repo       *repository.Repository
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

	LineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#FFFFFF"}).
			PaddingLeft(2)

	SelectedLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#22EE82"}).
				PaddingLeft(0)
	// ColumnStyle = lipgloss.NewStyle().Align(lipgloss.Left).Foreground(lipgloss.Color("#EBB2DF"))
	// MarginRight(2).
	// Height(2).
	// Width(2)
)
