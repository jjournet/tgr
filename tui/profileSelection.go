package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	m.OwnerList.SetHeight(h - headerHeight - footerHeight - 2 - 1)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2 - 5)
}

func InitProfileSelection() (tea.Model, tea.Cmd) {
	m := profileSelection{}
	m.InitTop("Profile Selection", constants.User.Login)
	m.InitBottom()
	columns := []table.Column{
		{Title: "Profile", Width: 50},
		{Title: "Description", Width: 50},
	}
	rows := []table.Row{}
	for _, owner := range constants.User.Owners {
		rows = append(rows, table.Row{owner.Login, owner.Description})
	}
	m.OwnerList = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

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
			owner := m.OwnerList.SelectedRow()[0]
			constants.Pr = profile.NewProfile(constants.User.Login, owner, constants.User.Client)
			return InitRepoSelection()
		}
	}
	m.OwnerList, _ = m.OwnerList.Update(msg)
	return m, nil
}

func (m profileSelection) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.OwnerList.View()), m.Bottom)
}
