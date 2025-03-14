package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lgtable "github.com/charmbracelet/lipgloss/table"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

type (
	errMsg struct{ error }
)

type item struct {
	Name string
}

func (i item) Title() string {
	return i.Name
}

func (i item) Description() string {
	return i.Name
}

func (i item) FilterValue() string {
	return i.Name
}

type ProfileModel struct {
	top      string
	repolist list.Model
	// bottom   string
	quitting bool
}

var (
	subtle      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	columnWidth = 40
	titleStyle  = func() lipgloss.Style {
		// b := lipgloss.NormalBorder()
		// b.TopLeft = "├"
		// return lipgloss.NewStyle().BorderStyle(b).Padding(0, 0)
		return lipgloss.NewStyle().Padding(0, 0).Border(lipgloss.NormalBorder(), true, true, false, true)
	}()

	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2).
			Height(8).
			Width(columnWidth + 1)

	listItemStyle = lipgloss.NewStyle().PaddingLeft(2).Render

	// mainStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#EBB2DF"))
	mainStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		b.TopLeft = "├"
		return lipgloss.NewStyle().BorderStyle(b)
	}()
)

func InitProfile() (tea.Model, tea.Cmd) {
	items := make([]list.Item, len(constants.Pr.Repos()))
	for i, repo := range constants.Pr.Repos() {
		items[i] = list.Item(item{Name: repo})
	}
	m := ProfileModel{repolist: list.New(items, list.NewDefaultDelegate(), 8, 8)}
	// m.list.Title = "Select a repository"
	// m.top = listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listItemStyle(fmt.Sprintf("%-20s %-30s", "Org:", constants.Pr.Org))))
	toprows := [][]string{
		{"Org", constants.Pr.Profile},
	}
	m.top = titleStyle.Render(lgtable.New().Rows(toprows...).Render())
	m.top = listStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listItemStyle(fmt.Sprintf("%-20s %-30s", "Org:", constants.Pr.Profile))))

	m.repolist.SetShowStatusBar(false)
	m.repolist.SetShowTitle(false)
	if constants.WindowSize.Height != 0 {
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		m.repolist.SetSize(constants.WindowSize.Width-2, constants.WindowSize.Height-1-headerHeight-footerHeight-2)
		mainStyle = mainStyle.Width(constants.WindowSize.Width - 2).Height(constants.WindowSize.Height - headerHeight - footerHeight - 2)
	}
	return m, func() tea.Msg { return errMsg{} }
}

func (m ProfileModel) Init() tea.Cmd {
	return nil
}

func (m ProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		m.repolist.SetSize(constants.WindowSize.Width-2, constants.WindowSize.Height-1-headerHeight-footerHeight-2)
		mainStyle = mainStyle.Width(constants.WindowSize.Width - 2).Height(constants.WindowSize.Height - headerHeight - footerHeight - 2)

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			// get the selected repo
			repoName := m.repolist.SelectedItem().(item).Name
			constants.Repo = repository.NewRepository(repoName, constants.Pr.Profile, constants.User.Client)
			return InitRepo()
		}
	}

	var cmd tea.Cmd
	m.repolist, cmd = m.repolist.Update(msg)
	return m, cmd
}

func (m ProfileModel) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.mainView(), m.footerView())
}

func (m ProfileModel) headerView() string {
	return titleStyle.Render(lipgloss.JoinVertical(lipgloss.Left, listItemStyle(fmt.Sprintf("%-20s %-30s", "Org:", constants.Pr.Profile))))
}

func (m ProfileModel) footerView() string {
	return "status bar"
}

func (m ProfileModel) mainView() string {
	return mainStyle.Render(m.repolist.View())
}
