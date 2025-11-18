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

type repoView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner    string
	repoName string

	// State
	repoDetails *github.RepoDetails
	workflows   []github.WorkflowInfo
	issues      []github.IssueInfo
	loading     bool
	err         error

	// UI
	EltList table.Model
}

func (m *repoView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewRepoView creates a new repository view model
func NewRepoView(ghService *github.GitHubService, owner, repoName string) (tea.Model, tea.Cmd) {
	m := &repoView{
		ghService: ghService,
		owner:     owner,
		repoName:  repoName,
		loading:   true,
	}

	m.InitTop(owner, repoName, "Loading...")
	m.TopFields = []string{owner, repoName, "Repository Summary"}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select", "(backspace) Back"}

	// Load repo details and workflows asynchronously
	return m, tea.Batch(
		ghService.LoadRepoDetailsCmd(owner, repoName),
		ghService.LoadWorkflowsCmd(owner, repoName),
		ghService.LoadIssuesCmd(owner, repoName),
	)
}

func (m *repoView) Init() tea.Cmd {
	return nil
}

func (m *repoView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.RepoDetailsLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		m.repoDetails = msg.Repo
		m.checkLoadingComplete()
		return m, nil

	case github.WorkflowsLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		m.workflows = msg.Workflows
		m.checkLoadingComplete()
		return m, nil

	case github.IssuesLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		m.issues = msg.Issues
		m.checkLoadingComplete()
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
			return NewRepoSelection(m.ghService, m.owner, true) // TODO: track isUser properly
		case "enter":
			// get the selected option
			row := m.EltList.HighlightedRow()
			if row.Data["id"] == types.WORKFLOW {
				return NewWorkflowList(m.ghService, m.owner, m.repoName)
			}
			if row.Data["id"] == types.ISSUE {
				return NewIssueList(m.ghService, m.owner, m.repoName)
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

func (m *repoView) checkLoadingComplete() {
	if m.repoDetails != nil && m.workflows != nil && m.issues != nil {
		m.loading = false
		m.EltList = m.buildSummaryListModel()

		if constants.WindowSize.Height != 0 {
			m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
		}
	}
}

func (m *repoView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading repository information..."
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

func (m *repoView) buildSummaryListModel() table.Model {
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
		"value":     m.repoName,
		"id":        types.PROJECT,
	}))

	// Display Description
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.DESCRIPTION),
		"value":     fmt.Sprintf("Description: %s", m.repoDetails.Description),
		"id":        types.DESCRIPTION,
	}))

	// Display Workflow
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.WORKFLOW),
		"value":     fmt.Sprintf("Workflows: %d", len(m.workflows)),
		"id":        types.WORKFLOW,
	}))

	// Display Issues
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      types.ConvertRepoElementType(types.ISSUE),
		"value":     fmt.Sprintf("Issues: %d", len(m.issues)),
		"id":        types.ISSUE,
	}))

	// Display Languages
	var langs string
	for lang, size := range m.repoDetails.Languages {
		langs += fmt.Sprintf("%s (%d) ", lang, size)
	}
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      "Languages",
		"value":     langs,
		"id":        types.LANGUAGES,
	}))

	return table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).
		WithHeaderVisibility(false).
		WithHighlightedRow(0)
}
