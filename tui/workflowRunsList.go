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
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("id", "ID", 10),
		table.NewColumn("title", "Title", 40),
		table.NewColumn("status", "Status", 10),
		table.NewColumn("created_at", "Created At", 20),
		table.NewColumn("branch", "Branch", 20),
	}
	rows := []table.Row{}
	for _, run := range constants.Runs {
		rows = append(rows, table.Row{
			Data: map[string]any{
				"arrow":      "",
				"id":         run.ID,
				"title":      run.Title,
				"status":     run.Status,
				"created_at": run.CreatedAt,
			},
		})
	}
	noBorder := table.Border{}
	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(noBorder).
		WithBaseStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#77c2f9")).Align(lipgloss.Left).Padding(0, 1)).
		HighlightStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#22EE82")).Background(lipgloss.Color("#111111")).Padding(0, 1)).
		Filtered(true)
}
