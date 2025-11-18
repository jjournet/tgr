package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui"
)

func main() {
	// Setup logging
	f, err := tea.LogToFile("log.txt", "debug")
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Starting tgr")

	// Create centralized GitHub service
	ghService, err := github.NewGitHubServiceFromCLI()
	if err != nil {
		log.Fatalf("Error creating GitHub service: %v", err)
	}

	// Create initial model
	initialModel := tui.NewApp(ghService)

	// Start the program
	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
