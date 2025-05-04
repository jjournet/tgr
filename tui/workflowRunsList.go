package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/tui/constants"
	"github.com/jjournet/tgr/types"
)

type workflowRunListView struct {
	commonElements
	EltList table.Model
}

func (m *workflowRunListView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

func InitWorkflowRunList() (tea.Model, tea.Cmd) {
	m := workflowRunListView{}
	m.InitTop("Workflow Run List", constants.Repo.GetRepoName())
	m.InitBottom()

	m.EltList = GetWorkflowRunListModel()

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m workflowRunListView) Init() tea.Cmd {
	return nil
}

func (m workflowRunListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			// get the selected option
			row := m.EltList.HighlightedRow()
			if row.Data["id"] == types.WORKFLOW {
				return InitWorflowList()
			}
		}
	}
	var cmd tea.Cmd
	return m, cmd
}

func (m workflowRunListView) View() string {
	for i, row := range m.EltList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.EltList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.EltList.View()), m.Bottom)
}

func GetWorkflowRunListModel() table.Model {
	columns := []table.Column{
    table.NewColumn("arrow", " ", table.WithWidth(3)),
    table.NewColumn("id", "ID", table.WithWidth(10)),
    table.NewColumn("title", "Title", table.WithWidth(20)),
    table.NewColumn("status", "Status", table.WithWidth(10)),
    table.NewColumn("created_at", "Created At", table.WithWidth(20)),
    
}
