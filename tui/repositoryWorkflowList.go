package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
	"github.com/jjournet/tgr/types"
)

type repoWorkflowListView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner    string
	repoName string

	// State
	workflows []github.WorkflowInfo
	loading   bool
	err       error

	// UI
	EltList table.Model
}

func (m *repoWorkflowListView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewWorkflowList creates a new workflow list view model
func NewWorkflowList(ghService *github.GitHubService, owner, repoName string) (tea.Model, tea.Cmd) {
	m := &repoWorkflowListView{
		ghService: ghService,
		owner:     owner,
		repoName:  repoName,
		loading:   true,
	}

	m.InitTop(owner, repoName, "Workflow List")
	m.TopFields = []string{owner, repoName, "Loading workflows..."}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) View Runs", "(t) Trigger", "(backspace) Back"}

	// Load workflows asynchronously
	return m, ghService.LoadWorkflowsCmd(owner, repoName)
}

func (m *repoWorkflowListView) Init() tea.Cmd {
	return nil
}

func (m *repoWorkflowListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.WorkflowsLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.workflows = msg.Workflows
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("Workflow List (%d workflows)", len(m.workflows))

		// Build UI table
		m.EltList = m.buildWorkflowListModel()

		if constants.WindowSize.Height != 0 {
			m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
		}

		return m, nil

	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "backspace":
			return NewRepoView(m.ghService, m.owner, m.repoName)
		case "enter":
			// get the selected option
			row := m.EltList.HighlightedRow()
			if row.Data["type"] == types.WORKFLOW {
				workflowID := row.Data["id"].(int64)
				return NewWorkflowRunList(m.ghService, m.owner, m.repoName, workflowID)
			}
		case "t", "T":
			// Trigger the selected workflow
			row := m.EltList.HighlightedRow()
			if row.Data["type"] == types.WORKFLOW {
				workflowID := row.Data["id"].(int64)
				return NewWorkflowInputForm(m.ghService, m.owner, m.repoName, workflowID, m), nil
			}
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.EltList, cmd = m.EltList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *repoWorkflowListView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading workflows..."
	}

	for i, row := range m.EltList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.EltList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.Render(m.EltList.View()),
		m.RenderBottomFields(),
	)
}

func (m *repoWorkflowListView) buildWorkflowListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("workflow", "Workflow", 30),
		table.NewColumn("state", "State", 100),
		table.NewColumn("id", "ID", 10),
		table.NewColumn("type", "Type", 10),
	}

	rows := []table.Row{}
	for _, workflow := range m.workflows {
		rows = append(rows, table.NewRow(table.RowData{
			"workflow": workflow.Name,
			"state":    workflow.State,
			"id":       workflow.ID,
			"type":     types.WORKFLOW,
		}))
	}

	noBorder := table.Border{}
	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(noBorder).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true)
}
