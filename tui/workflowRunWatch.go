package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type workflowRunWatchView struct {
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
	jobs      []github.JobInfo
	loading   bool
	err       error
	viewport  viewport.Model

	// Refresh
	refreshInterval time.Duration
}

type tickMsg time.Time

func (m *workflowRunWatchView) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
	m.viewport.Width = w - 4 // Account for borders/padding
	m.viewport.Height = h - headerHeight - footerHeight - 2
}

// NewWorkflowRunWatch creates a new workflow run watch view model
func NewWorkflowRunWatch(ghService *github.GitHubService, owner, repoName string, workflowID, runID int64) (tea.Model, tea.Cmd) {
	m := &workflowRunWatchView{
		ghService:       ghService,
		owner:           owner,
		repoName:        repoName,
		workflowID:      workflowID,
		runID:           runID,
		loading:         true,
		refreshInterval: 5 * time.Second,
		viewport:        viewport.New(0, 0),
	}

	m.InitTop(owner, repoName, fmt.Sprintf("Watching run #%d...", runID))
	m.TopFields = []string{owner, repoName, fmt.Sprintf("Watch Run #%d", runID)}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(backspace) Back", "(r) Refresh Now"}

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}

	// Initial load
	return m, tea.Batch(
		ghService.LoadRunDetailCmd(owner, repoName, runID),
		ghService.LoadRunJobsCmd(owner, repoName, runID),
		m.tick(),
	)
}

func (m *workflowRunWatchView) tick() tea.Cmd {
	return tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *workflowRunWatchView) Init() tea.Cmd {
	return nil
}

func (m *workflowRunWatchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		// Only refresh if we are still viewing this page
		// In a real app we might want to check if the run is completed to stop refreshing
		if m.runDetail != nil && m.runDetail.Status == "completed" {
			return m, nil
		}
		return m, tea.Batch(
			m.ghService.LoadRunDetailCmd(m.owner, m.repoName, m.runID),
			m.ghService.LoadRunJobsCmd(m.owner, m.repoName, m.runID),
			m.tick(),
		)

	case github.RunDetailLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.runDetail = msg.Run
		m.checkLoadingComplete()
		atBottom := m.viewport.AtBottom()
		m.viewport.SetContent(m.renderContent())
		if atBottom {
			m.viewport.GotoBottom()
		}
		return m, nil

	case github.RunJobsLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.jobs = msg.Jobs
		m.checkLoadingComplete()
		atBottom := m.viewport.AtBottom()
		m.viewport.SetContent(m.renderContent())
		if atBottom {
			m.viewport.GotoBottom()
		}
		return m, nil

	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "backspace":
			return NewWorkflowRunList(m.ghService, m.owner, m.repoName, m.workflowID)
		case "r":
			return m, tea.Batch(
				m.ghService.LoadRunDetailCmd(m.owner, m.repoName, m.runID),
				m.ghService.LoadRunJobsCmd(m.owner, m.repoName, m.runID),
			)
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *workflowRunWatchView) checkLoadingComplete() {
	if m.runDetail != nil && m.jobs != nil {
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("Watch Run #%d - %s (%s)", m.runDetail.RunNumber, m.runDetail.Name, m.runDetail.Status)
	}
}

func (m *workflowRunWatchView) renderContent() string {
	var content strings.Builder

	// Run Status Header
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

	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)
	content.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, strings.ToUpper(m.runDetail.Status))))
	if m.runDetail.Conclusion != "" {
		content.WriteString(statusStyle.Render(fmt.Sprintf(" - %s", strings.ToUpper(m.runDetail.Conclusion))))
	}

	// Elapsed time
	duration := time.Since(m.runDetail.CreatedAt)
	if !m.runDetail.UpdatedAt.IsZero() && m.runDetail.Status == "completed" {
		duration = m.runDetail.UpdatedAt.Sub(m.runDetail.CreatedAt)
	}
	content.WriteString(fmt.Sprintf("  (%s)", duration.Round(time.Second)))
	content.WriteString("\n\n")

	// Jobs and Steps
	for _, job := range m.jobs {
		jobIcon := "○"
		jobColor := lipgloss.Color("#8B949E") // Grey

		if job.Status == "completed" && job.Conclusion == "success" {
			jobIcon = "✓"
			jobColor = lipgloss.Color("#22EE82")
		} else if job.Status == "in_progress" {
			jobIcon = "●" // or spinner
			jobColor = lipgloss.Color("#FFCC00")
		} else if job.Status == "queued" {
			jobIcon = "○"
			jobColor = lipgloss.Color("#FFCC00")
		} else if job.Status == "completed" && job.Conclusion == "failure" {
			jobIcon = "✗"
			jobColor = lipgloss.Color("#FF0000")
		}

		jobStyle := lipgloss.NewStyle().Foreground(jobColor).Bold(true)
		content.WriteString(jobStyle.Render(fmt.Sprintf("%s %s", jobIcon, job.Name)))

		// Job duration
		if !job.StartedAt.IsZero() {
			jobEnd := time.Now()
			if !job.CompletedAt.IsZero() {
				jobEnd = job.CompletedAt
			}
			jobDuration := jobEnd.Sub(job.StartedAt)
			content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#5865F2")).Render(fmt.Sprintf(" (%s)", jobDuration.Round(time.Second))))
		}
		content.WriteString("\n")

		// Steps
		for _, step := range job.Steps {
			stepIcon := "  ○"
			stepColor := lipgloss.Color("#8B949E")

			if step.Status == "completed" && step.Conclusion == "success" {
				stepIcon = "  ✓"
				stepColor = lipgloss.Color("#22EE82")
			} else if step.Status == "in_progress" {
				stepIcon = "  ●"
				stepColor = lipgloss.Color("#FFCC00")
			} else if step.Status == "queued" {
				stepIcon = "  ○"
				stepColor = lipgloss.Color("#8B949E")
			} else if step.Status == "completed" && step.Conclusion == "failure" {
				stepIcon = "  ✗"
				stepColor = lipgloss.Color("#FF0000")
			} else if step.Status == "completed" && step.Conclusion == "skipped" {
				stepIcon = "  -"
				stepColor = lipgloss.Color("#8B949E")
			}

			stepStyle := lipgloss.NewStyle().Foreground(stepColor)
			content.WriteString(stepStyle.Render(fmt.Sprintf("%s %s", stepIcon, step.Name)))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}
	return content.String()
}

func (m *workflowRunWatchView) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading run details and jobs..."
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.Render(m.viewport.View()),
		m.RenderBottomFields(),
	)
}
