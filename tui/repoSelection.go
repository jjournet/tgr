package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

type repoSelection struct {
	commonElements
	RepoList table.Model
}

func (m *repoSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	m.RepoList.SetHeight(h - headerHeight - footerHeight - 2 - 1)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2 - 5)
	log.Println("new height: ", h, "new width: ", w)
	log.Println("Repo set height: ", h-headerHeight-footerHeight-2-1, "Main set height: ", h-headerHeight-footerHeight-2-5)
}

func InitRepoSelection() (tea.Model, tea.Cmd) {
	m := repoSelection{}
	m.InitTop("Repo Selection", constants.Pr.Profile)
	m.InitBottom()
	columns := []table.Column{
		{Title: "Repo", Width: 50},
		{Title: "Description", Width: 50},
	}
	rows := []table.Row{}
	for _, repo := range constants.Pr.RepoList {
		rows = append(rows, table.Row{repo.Name, repo.Description})
	}
	m.RepoList = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m repoSelection) Init() tea.Cmd {
	return nil
}

func (m repoSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			// get the selected repo
			repoName := m.RepoList.SelectedRow()[0]
			constants.Repo = repository.NewRepository(repoName, constants.Pr.Profile, constants.User.Client)
			return InitRepo()
		}
	}
	m.RepoList, _ = m.RepoList.Update(msg)
	return m, nil
}

func (m repoSelection) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.RepoList.View()), m.Bottom)
}
