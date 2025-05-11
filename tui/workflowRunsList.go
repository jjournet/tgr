package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/runs"
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

func InitWorkflowRunList(wkflwId int64) (tea.Model, tea.Cmd) {
	m := workflowRunListView{}
	m.InitTop(constants.Pr.Profile, constants.Repo.GetRepoName(), fmt.Sprintf("Workflow Run List for %d", wkflwId))
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
			return InitWorflowList()
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
		table.NewColumn("indicator", " ", 3),
		table.NewColumn("title", "Title", 35),
		table.NewColumn("status", "Status", 15),
		table.NewColumn("conclusion", "Conclusion", 15),
		table.NewColumn("created_at", "Created At", 24),
		table.NewColumn("branch", "Branch", 20),
		table.NewColumn("id", "ID", 12),
	}
	rows := []table.Row{}

	for _, run := range constants.Runs {
		rows = append(rows, makeRow(*run))
	}

	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).
		WithHighlightedRow(0).
		WithFooterVisibility(false)

}

func makeRow(run runs.Run) table.Row {
	indicator := "\uf128"
	color := lipgloss.Color("#77c2f9")
	if run.Status == "completed" && run.Conclusion == "success" {
		indicator = "\uf058"
		color = lipgloss.Color("#22EE82")
	} else if run.Status == "in_progress" {
		indicator = "\uef0c"
		color = lipgloss.Color("#FFCC00")
	} else if run.Status == "queued" {
		indicator = "󰚭"
		color = lipgloss.Color("#FFCC00")
	} else if run.Status == "completed" && run.Conclusion == "failure" {
		indicator = "\uea87"
		color = lipgloss.Color("#FF0000")
	}
	// return table.Row{
	// 	Data: map[string]any{
	// 		"arrow":      "",
	// 		"indicator":  indicator,
	// 		"title":      run.Title,
	// 		"status":     run.Status,
	// 		"created_at": run.CreatedAt[0:19],
	// 		"branch":     run.Branch,
	// 		"id":         run.ID,
	// 	},
	// }

	return table.NewRow(table.RowData{
		"arrow":      "",
		"indicator":  table.NewStyledCell(indicator, lipgloss.NewStyle().Foreground(lipgloss.Color(color))),
		"title":      run.Title,
		"status":     run.Status,
		"conclusion": run.Conclusion,
		"created_at": run.CreatedAt[0:19],
		"branch":     run.Branch,
		"id":         run.ID,
	})
	// return table.NewRow(table.RowData{"arrow": "",
}
