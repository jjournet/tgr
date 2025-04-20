package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/tui/constants"
	"github.com/jjournet/tgr/types"
)

type repoWorkflowListView struct {
	commonElements
	EltList table.Model
}

func (m *repoWorkflowListView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

func InitWorflowList() (tea.Model, tea.Cmd) {
	m := repoView{}
	m.InitTop("Workflow List", constants.Repo.GetRepoName())
	m.InitBottom()

	m.EltList = GetWorkflowListModel()

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m repoWorkflowListView) Init() tea.Cmd {
	return nil
}

func (m repoWorkflowListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	m.EltList, cmd = m.EltList.Update(msg)
	return m, cmd
}

func (m repoWorkflowListView) View() string {
	for i, row := range m.EltList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.EltList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.EltList.View()), m.Bottom)
}
func GetWorkflowListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("workflow", "Workflow", 30),
		table.NewColumn("state", "State", 100),
		table.NewColumn("id", "ID", 10),
	}
	rows := []table.Row{}
	for _, workflow := range constants.Repo.GetWorkflows() {
		rows = append(rows, table.NewRow(table.RowData{"workflow": workflow.Name, "state": workflow.State, "id": types.WORKFLOW}))
	}
	noBorder := table.Border{}
	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(noBorder).
		WithBaseStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#77c2f9")).Align(lipgloss.Left).Padding(0, 1)).
		HighlightStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#22EE82")).Background(lipgloss.Color("#111111")).Padding(0, 1)).
		Filtered(true)
}
