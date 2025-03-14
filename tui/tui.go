package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/ghuser"
	"github.com/jjournet/tgr/tui/constants"
)

// StartTea starts the bubbletea program
// func StartTea(pr *profile.Profile) {
// 	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
// 		fmt.Println("Couldn't open a file for logging:", err)
// 		os.Exit(1)
// 	} else {
// 		defer func() {
// 			err = f.Close()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}()
// 	}

// 	constants.Pr = pr
// 	m, _ := InitProfile()
// 	constants.P = tea.NewProgram(m, tea.WithAltScreen())
// 	if err := constants.P.Start(); err != nil {
// 		fmt.Println("Error starting program:", err)
// 		os.Exit(1)
// 	}
// }

func StartTea(ghuser *ghuser.GHUser) {
	// if f, err := tea.LogToFile("debug.log", "help"); err != nil {
	// 	fmt.Println("Couldn't open a file for logging:", err)
	// 	os.Exit(1)
	// } else {
	// 	defer func() {
	// 		err = f.Close()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}()
	// }
	constants.User = ghuser
	m, _ := InitProfileSelection()
	// m, _ := InitOrgs()
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if err := constants.P.Start(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
