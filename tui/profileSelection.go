package tui

import (
	"fmt"

	// "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/profile"
	"github.com/jjournet/tgr/tui/constants"
)

type profileSelection struct {
	commonElements
	OwnerList table.Model
}

func (m *profileSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

func InitProfileSelection() (tea.Model, tea.Cmd) {
	m := profileSelection{}
	m.InitTop("Profile Selection", constants.User.Login)
	m.InitBottom()

	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("profile", "Profile", 30),
		table.NewColumn("desc", "Description", 20),
	}
	rows := []table.Row{}
	for _, owner := range constants.User.Owners {
		rows = append(rows, table.NewRow(table.RowData{"profile": owner.Login, "desc": owner.Description}))
	}
	noBorder := table.Border{}
	m.OwnerList = table.New(columns).WithRows(rows).
		Focused(true).
		Border(noBorder).
		WithBaseStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#77c2f9")).Align(lipgloss.Left).Padding(0, 1)).
		HighlightStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#22EE82")).Background(lipgloss.Color("#111111")).Padding(0, 1)).
		Filtered(true)
	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m profileSelection) Init() tea.Cmd {
	return nil
}

func (m profileSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			owner := m.OwnerList.HighlightedRow().Data["profile"].(string)
			constants.Pr = profile.NewProfile(constants.User.Login, owner, constants.User.Client)
			return InitRepoSelection()
		}
	}

	m.OwnerList, _ = m.OwnerList.Update(msg)
	return m, nil
}

func (m profileSelection) View() string {
	for i, row := range m.OwnerList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.OwnerList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}

	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.OwnerList.View()), m.Bottom)
}
