package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type workflowRunDetailView struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner      string
	repoName   string
	workflowID int64
	runID      int64

	// State
	runDetail *github.RunDetailInfo
	loading   bool
	err       error
}

func (m *workflowRunDetailView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewWorkflowRunDetail creates a new workflow run detail view model
func NewWorkflowRunDetail(ghService *github.GitHubService, owner, repoName string, workflowID, runID int64) (tea.Model, tea.Cmd) {
	m := &workflowRunDetailView{
		ghService:  ghService,
		owner:      owner,
		repoName:   repoName,
		workflowID: workflowID,
		runID:      runID,
		loading:    true,
	}

	m.InitTop(owner, repoName, fmt.Sprintf("Loading run #%d...", runID))
	m.TopFields = []string{owner, repoName, fmt.Sprintf("Run #%d", runID)}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(backspace) Back"}

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}

	return m, ghService.LoadRunDetailCmd(owner, repoName, runID)
}

func (m *workflowRunDetailView) Init() tea.Cmd {
	return nil
}

func (m *workflowRunDetailView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.RunDetailLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.runDetail = msg.Run
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("Run #%d - %s", m.runDetail.RunNumber, m.runDetail.Name)

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
			return NewWorkflowRunList(m.ghService, m.owner, m.repoName, m.workflowID)
		}
	}

	return m, nil
}

func (m *workflowRunDetailView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading run details..."
	}

	var content strings.Builder

	// Status and Conclusion
	statusColor := lipgloss.Color("#77c2f9")
	statusIcon := "\uf128"

	if m.runDetail.Status == "completed" && m.runDetail.Conclusion == "success" {
		statusIcon = "\uf058"
		statusColor = lipgloss.Color("#22EE82")
	} else if m.runDetail.Status == "in_progress" {
		statusIcon = "\uef0c"
		statusColor = lipgloss.Color("#FFCC00")
	} else if m.runDetail.Status == "queued" {
		statusIcon = "󰚭"
		statusColor = lipgloss.Color("#FFCC00")
	} else if m.runDetail.Status == "completed" && m.runDetail.Conclusion == "failure" {
		statusIcon = "\uea87"
		statusColor = lipgloss.Color("#FF0000")
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	content.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, strings.ToUpper(m.runDetail.Status))))
	if m.runDetail.Conclusion != "" {
		content.WriteString(statusStyle.Render(fmt.Sprintf(" - %s", strings.ToUpper(m.runDetail.Conclusion))))
	}
	content.WriteString("\n\n")

	content.WriteString(titleStyle.Render(m.runDetail.Name))
	content.WriteString("\n\n")

	// Metadata
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B949E")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C9D1D9"))

	content.WriteString(labelStyle.Render("Run Number: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("#%d", m.runDetail.RunNumber)))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Attempt: "))
	content.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.runDetail.RunAttempt)))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Branch: "))
	content.WriteString(valueStyle.Render(m.runDetail.Branch))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Event: "))
	content.WriteString(valueStyle.Render(m.runDetail.Event))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Actor: "))
	content.WriteString(valueStyle.Render(m.runDetail.Actor))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Commit SHA: "))
	content.WriteString(valueStyle.Render(m.runDetail.HeadSHA[:8]))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Created: "))
	content.WriteString(valueStyle.Render(m.runDetail.CreatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Updated: "))
	content.WriteString(valueStyle.Render(m.runDetail.UpdatedAt.Format("2006-01-02 15:04:05")))
	content.WriteString("\n\n")

	// URLs
	content.WriteString(strings.Repeat("─", constants.WindowSize.Width-4))
	content.WriteString("\n\n")

	urlStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#58A6FF")).
		Underline(true)

	content.WriteString(labelStyle.Render("GitHub URL: "))
	content.WriteString(urlStyle.Render(m.runDetail.HTMLURL))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Jobs URL: "))
	content.WriteString(urlStyle.Render(m.runDetail.JobsURL))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Logs URL: "))
	content.WriteString(urlStyle.Render(m.runDetail.LogsURL))
	content.WriteString("\n")

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.Render(content.String()),
		m.RenderBottomFields(),
	)
}
