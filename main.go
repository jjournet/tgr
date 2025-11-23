package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/config"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui"
)

func main() {
	// Parse flags
	loginFlag := flag.Bool("login", false, "Force login window to update credentials")
	flag.Parse()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to default if loading fails, but try to log it later
		cfg = &config.Config{LogLevel: "INFO"}
	}

	// Setup logging
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()

	// Set log level
	var level slog.Level
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(f, opts))
	slog.SetDefault(logger)
	// Redirect standard log to file as well for dependencies
	log.SetOutput(f)

	slog.Debug("Starting tgr")

	// Initialize Auth Service
	authService, err := github.NewAuthService()
	if err != nil {
		slog.Error("Error initializing auth service", "error", err)
		os.Exit(1)
	}

	// Check credentials
	username, token, err := authService.GetCredentials()

	// Determine if we need to show login
	showLogin := *loginFlag || err != nil || username == "" || token == ""

	if showLogin {
		// Run Login Program
		loginModel := tui.NewLoginModel(authService, username, token)
		p := tea.NewProgram(loginModel)
		m, err := p.Run()
		if err != nil {
			slog.Error("Error running login", "error", err)
			os.Exit(1)
		}

		// Check if login was successful
		finalLoginModel, ok := m.(tui.LoginModel)
		if !ok {
			// Should not happen
			os.Exit(1)
		}

		// If user quit without submitting
		username, token = finalLoginModel.GetCredentials()
		// We need to verify if they actually submitted successfully or just quit
		// The LoginModel should probably expose a "Success" field or we check credentials again

		// Re-fetch to be sure
		username, token, err = authService.GetCredentials()
		if err != nil || username == "" || token == "" {
			slog.Debug("Login cancelled or failed")
			os.Exit(0)
		}
	}

	// Create centralized GitHub service
	ghService := github.NewGitHubService(token)

	// Create initial model
	initialModel := tui.NewApp(ghService)

	// Start the program
	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}
}
