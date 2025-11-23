package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/github"
)

// App is the root model that manages the application
type App struct {
	ghService   *github.GitHubService
	currentView tea.Model
	err         error
	initCmd     tea.Cmd // Store initial command to run in Init()
}

// NewApp creates the root application model
func NewApp(ghService *github.GitHubService) *App {
	// Start with profile selection
	profileView, initCmd := NewProfileSelection(ghService)

	// Store the initial command to be returned from Init()
	return &App{
		ghService:   ghService,
		currentView: profileView,
		initCmd:     initCmd,
	}
}

func (a *App) Init() tea.Cmd {
	slog.Debug("App.Init() called")
	// Return the initial command from the first view
	if a.initCmd != nil {
		slog.Debug("App.Init() returning stored initCmd")
		cmd := a.initCmd
		a.initCmd = nil // Clear it so it doesn't run again
		return cmd
	}
	slog.Debug("App.Init() calling currentView.Init()")
	return a.currentView.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Log all messages for debugging
	slog.Debug("App.Update received message", "type", fmt.Sprintf("%T", msg))

	// Handle global keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
	}

	// Forward to current view
	var cmd tea.Cmd
	a.currentView, cmd = a.currentView.Update(msg)

	return a, cmd
}

func (a *App) View() string {
	if a.err != nil {
		return a.err.Error()
	}
	return a.currentView.View()
}
