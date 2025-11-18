package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type issueDetailView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner    string
	repoName string
	issue    github.IssueInfo
}

func (m *issueDetailView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewIssueDetail creates a new issue detail view model
func NewIssueDetail(ghService *github.GitHubService, owner, repoName string, issue github.IssueInfo) (tea.Model, tea.Cmd) {
	m := &issueDetailView{
		ghService: ghService,
		owner:     owner,
		repoName:  repoName,
		issue:     issue,
	}

	m.InitTop(owner, repoName, fmt.Sprintf("Issue #%d", issue.Number))
	m.TopFields = []string{owner, repoName, fmt.Sprintf("Issue #%d", issue.Number)}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(backspace) Back"}

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}

	return m, nil
}

func (m *issueDetailView) Init() tea.Cmd {
	return nil
}

func (m *issueDetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "backspace":
			return NewIssueList(m.ghService, m.owner, m.repoName)
		}
	}

	return m, nil
}

func (m *issueDetailView) View() string {
	var content strings.Builder

	// State indicator and title
	stateColor := lipgloss.Color("#22EE82")
	stateIcon := "\uf468"
	stateText := "OPEN"
	if m.issue.State == "closed" {
		stateColor = lipgloss.Color("#b19cd9")
		stateIcon = "\uf46a"
		stateText = "CLOSED"
	}

	stateStyle := lipgloss.NewStyle().
		Foreground(stateColor).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	content.WriteString(stateStyle.Render(fmt.Sprintf("%s %s", stateIcon, stateText)))
	content.WriteString("  ")
	content.WriteString(titleStyle.Render(m.issue.Title))
	content.WriteString("\n\n")

	// Metadata
	metadataStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B949E"))

	content.WriteString(metadataStyle.Render(fmt.Sprintf("Author: %s", m.issue.Author)))
	content.WriteString("\n")
	content.WriteString(metadataStyle.Render(fmt.Sprintf("Created: %s", m.issue.CreatedAt.Format("2006-01-02 15:04:05"))))
	content.WriteString("\n")
	content.WriteString(metadataStyle.Render(fmt.Sprintf("Updated: %s", m.issue.UpdatedAt.Format("2006-01-02 15:04:05"))))
	content.WriteString("\n")
	content.WriteString(metadataStyle.Render(fmt.Sprintf("Comments: %d", m.issue.Comments)))
	content.WriteString("\n")

	// Labels
	if len(m.issue.Labels) > 0 {
		content.WriteString("\n")
		content.WriteString(metadataStyle.Render("Labels: "))
		labelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58A6FF")).
			Bold(true)
		content.WriteString(labelStyle.Render(strings.Join(m.issue.Labels, ", ")))
		content.WriteString("\n")
	}

	// Separator
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", constants.WindowSize.Width-4))
	content.WriteString("\n\n")

	// Body
	if m.issue.Body != "" {
		bodyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C9D1D9"))
		content.WriteString(bodyStyle.Render(m.issue.Body))
	} else {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B949E")).
			Italic(true)
		content.WriteString(emptyStyle.Render("No description provided."))
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.Render(content.String()),
		m.RenderBottomFields(),
	)
}
