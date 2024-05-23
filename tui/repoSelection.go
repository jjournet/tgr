package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

type repoSelection struct {
	commonElements
	// RepoList table.Model
	RepoList table.Model
}

func (m *repoSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

func InitRepoSelection() (tea.Model, tea.Cmd) {
	m := repoSelection{}
	m.InitTop("Repo Selection", constants.Pr.Profile)
	m.InitBottom()

	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("repo", "Repository", 40).WithFiltered(true),
		table.NewColumn("desc", "Description", 100),
	}
	items := []table.Row{}
	for _, repo := range constants.Pr.RepoList {
		items = append(items, table.NewRow(table.RowData{"repo": repo.Name, "desc": repo.Description}))
	}
	m.RepoList = table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true)

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
			repoName := m.RepoList.HighlightedRow().Data["repo"].(string)
			constants.Repo = repository.NewRepository(repoName, constants.Pr.Profile, constants.User.Client)
			return InitRepoView()
		case "backspace":
			return InitProfileSelection()
		}
	}
	var cmd tea.Cmd
	m.RepoList, cmd = m.RepoList.Update(msg)
	return m, cmd
}

func (m repoSelection) View() string {
	for i, row := range m.RepoList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.RepoList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.RepoList.View()), m.Bottom)
}
