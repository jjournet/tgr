package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/profile"
	"github.com/jjournet/tgr/tui/constants"
)

type orgItem struct {
	Name string
	Id   int64
}

func (i orgItem) Title() string {
	return i.Name
}

func (i orgItem) Description() string {
	return fmt.Sprintf("%d", i.Id)
}

func (i orgItem) FilterValue() string {
	return i.Name
}

type OrgModel struct {
	OrgList  list.Model
	quitting bool
}

func InitOrgs() (tea.Model, tea.Cmd) {
	items := make([]list.Item, len(constants.User.Orgs))
	for i, org := range constants.User.Orgs {
		items[i] = list.Item(orgItem{Name: org, Id: int64(i)})
	}
	m := OrgModel{OrgList: list.New(items, list.NewDefaultDelegate(), 8, 8)}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.OrgList.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	return m, func() tea.Msg { return errMsg{} }
}

func (m OrgModel) Init() tea.Cmd {
	return nil
}

func (m OrgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.OrgList.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// get the selected org
			org := m.OrgList.SelectedItem().(orgItem)
			constants.Pr = profile.NewProfile(constants.User.Login, org.Name, constants.User.Client)
			return InitProfile()
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			return InitProfile()
		}
	}
	var cmd tea.Cmd
	m.OrgList, cmd = m.OrgList.Update(msg)
	return m, cmd
}

func (m OrgModel) View() string {
	return m.OrgList.View()
}
