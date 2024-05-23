package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/tui/constants"
	repositoryviews "github.com/jjournet/tgr/tui/repositoryViews"
)

type repoView struct {
	commonElements
	EltList table.Model
}

func (m *repoView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

func InitRepoView() (tea.Model, tea.Cmd) {
	m := repoView{}
	m.InitTop("Repository Summary", constants.Repo.Name)
	m.InitBottom()

	m.EltList = repositoryviews.GetSummaryListModel()

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m repoView) Init() tea.Cmd {
	return nil
}

func (m repoView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "backspace":
			return InitRepoSelection()
		}
	}
	var cmd tea.Cmd
	m.EltList, cmd = m.EltList.Update(msg)
	return m, cmd
}

func (m repoView) View() string {
	for i, row := range m.EltList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.EltList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.EltList.View()), m.Bottom)
}
