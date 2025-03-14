package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/jjournet/tgr/ghuser"
	"github.com/jjournet/tgr/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v69/github"
)

// var token string = ""

var client *github.Client = nil

func initGH() {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}
	token := strings.TrimSuffix(string(out), "\n")
	client = github.NewClient(nil).WithAuthToken(token)

}

func main() {
	f, err := tea.LogToFile("log.txt", "debug")
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting tgr")
	initGH()
	// pr := profile.NewProfile("jjournet_HQY01", client)
	user := ghuser.NewUser(client)
	tui.StartTea(user)
}
