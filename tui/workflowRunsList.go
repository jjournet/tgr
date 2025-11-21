package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type workflowRunListView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner      string
	repoName   string
	workflowID int64

	// State
	runs    []github.RunInfo
	loading bool
	err     error

	// UI
	EltList table.Model
}

func (m *workflowRunListView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewWorkflowRunList creates a new workflow run list view model
func NewWorkflowRunList(ghService *github.GitHubService, owner, repoName string, workflowID int64) (tea.Model, tea.Cmd) {
	m := &workflowRunListView{
		ghService:  ghService,
		owner:      owner,
		repoName:   repoName,
		workflowID: workflowID,
		loading:    true,
	}

	m.InitTop(owner, repoName, fmt.Sprintf("Loading runs for workflow %d...", workflowID))
	m.TopFields = []string{owner, repoName, fmt.Sprintf("Workflow Run List for %d", workflowID)}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select", "(w) Watch", "(backspace) Back"}

	// Load workflow runs asynchronously
	return m, ghService.LoadWorkflowRunsCmd(owner, repoName, workflowID)
}

func (m *workflowRunListView) Init() tea.Cmd {
	return nil
}

func (m *workflowRunListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.WorkflowRunsLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.runs = msg.Runs
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("Workflow Run List (%d runs)", len(m.runs))

		// Build UI table
		m.EltList = m.buildWorkflowRunListModel()

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
			return NewWorkflowList(m.ghService, m.owner, m.repoName)
		case "enter":
			// Get the selected run
			row := m.EltList.HighlightedRow()
			runID := row.Data["id"].(int64)
			return NewWorkflowRunDetail(m.ghService, m.owner, m.repoName, m.workflowID, runID)
		case "w":
			// Get the selected run
			row := m.EltList.HighlightedRow()
			runID := row.Data["id"].(int64)
			return NewWorkflowRunWatch(m.ghService, m.owner, m.repoName, m.workflowID, runID)
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.EltList, cmd = m.EltList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *workflowRunListView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading workflow runs..."
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

func (m *workflowRunListView) buildWorkflowRunListModel() table.Model {
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
	for _, run := range m.runs {
		rows = append(rows, makeRunRow(run))
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

func makeRunRow(run github.RunInfo) table.Row {
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

	createdAt := run.CreatedAt.Format("2006-01-02 15:04:05")

	return table.NewRow(table.RowData{
		"arrow":      "",
		"indicator":  table.NewStyledCell(indicator, lipgloss.NewStyle().Foreground(color)),
		"title":      run.Title,
		"status":     run.Status,
		"conclusion": run.Conclusion,
		"created_at": createdAt,
		"branch":     run.Branch,
		"id":         run.ID,
	})
}
