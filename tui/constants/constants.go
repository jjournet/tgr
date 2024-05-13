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
)

var (
	DocStyle = lipgloss.NewStyle().Margin(1)
)
