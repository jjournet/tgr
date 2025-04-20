package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/tui/constants"
	"github.com/jjournet/tgr/types"
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
	m.InitTop("Repository Summary", constants.Repo.GetRepoName())
	m.InitBottom()

	m.EltList = GetSummaryListModel()

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

func (m repoView) View() string {
	for i, row := range m.EltList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.EltList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.EltList.View()), m.Bottom)
}

func GetSummaryListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("indicator", " ", 3),
		table.NewColumn("type", "Repository info", 40).WithFiltered(true),
		table.NewColumn("value", "Value", 80),
	}

	var items []table.Row
	// Display Project
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.PROJECT),
		"value":     constants.Repo.GetRepoName(),
		"id":        types.PROJECT,
	}))
	// Display Description
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.DESCRIPTION),
		"value":     fmt.Sprintf("Description: %s", constants.Repo.GetDescription()),
		"id":        types.DESCRIPTION,
	}))
	// Display Workflow
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.WORKFLOW),
		"value":     fmt.Sprintf("Workflow: %d", len(constants.Repo.GetWorkflows())),
		"id":        types.WORKFLOW,
	}))
	items = append(items, table.NewRow(table.RowData{"indicator": "",
		"type":  types.ConvertRepoElementType(types.RUN),
		"value": fmt.Sprintf("Actions: %d", len(constants.Repo.GetRuns())),
		"id":    types.RUN,
	}))
	// append all languages in one string, with percentage in parenthesis
	var langs string
	languages := constants.Repo.GetLanguages()
	for lang := range languages {
		langs += fmt.Sprintf("%s (%d) ", lang, languages[lang])
	}
	items = append(items, table.NewRow(table.RowData{"indicator": "",
		"type":  "Languages",
		"value": langs,
		"id":    types.LANGUAGES,
	}))

	return table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).WithHeaderVisibility(false).
		WithHighlightedRow(0)

}
