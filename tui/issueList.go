package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type issueListView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner    string
	repoName string

	// State
	issues  []github.IssueInfo
	loading bool
	err     error

	// UI
	EltList table.Model
}

func (m *issueListView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewIssueList creates a new issue list view model
func NewIssueList(ghService *github.GitHubService, owner, repoName string) (tea.Model, tea.Cmd) {
	m := &issueListView{
		ghService: ghService,
		owner:     owner,
		repoName:  repoName,
		loading:   true,
	}

	m.InitTop(owner, repoName, "Loading issues...")
	m.TopFields = []string{owner, repoName, "Issue List"}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select", "(backspace) Back"}

	// Load issues asynchronously
	return m, ghService.LoadIssuesCmd(owner, repoName)
}

func (m *issueListView) Init() tea.Cmd {
	return nil
}

func (m *issueListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.IssuesLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.issues = msg.Issues
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("Issue List (%d issues)", len(m.issues))

		// Build UI table
		m.EltList = m.buildIssueListModel()

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
			// Get the selected issue
			row := m.EltList.HighlightedRow()
			issueNumber := row.Data["number"].(int)

			// Find the issue in our list
			for _, issue := range m.issues {
				if issue.Number == issueNumber {
					return NewIssueDetail(m.ghService, m.owner, m.repoName, issue)
				}
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

func (m *issueListView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading issues..."
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

func (m *issueListView) buildIssueListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("indicator", " ", 3),
		table.NewColumn("number", "#", 6),
		table.NewColumn("title", "Title", 50).WithFiltered(true),
		table.NewColumn("state", "State", 10),
		table.NewColumn("author", "Author", 20).WithFiltered(true),
		table.NewColumn("comments", "Comments", 10),
		table.NewColumn("labels", "Labels", 30).WithFiltered(true),
	}

	rows := []table.Row{}
	for _, issue := range m.issues {
		rows = append(rows, makeIssueRow(issue))
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

func makeIssueRow(issue github.IssueInfo) table.Row {
	indicator := "\uf128"
	color := lipgloss.Color("#77c2f9")

	if issue.State == "open" {
		indicator = "\uf468"
		color = lipgloss.Color("#22EE82")
	} else if issue.State == "closed" {
		indicator = "\uf46a"
		color = lipgloss.Color("#b19cd9")
	}

	labels := strings.Join(issue.Labels, ", ")
	if len(labels) > 30 {
		labels = labels[:27] + "..."
	}

	return table.NewRow(table.RowData{
		"arrow":     "",
		"indicator": table.NewStyledCell(indicator, lipgloss.NewStyle().Foreground(color)),
		"number":    issue.Number,
		"title":     issue.Title,
		"state":     issue.State,
		"author":    issue.Author,
		"comments":  issue.Comments,
		"labels":    labels,
	})
}
