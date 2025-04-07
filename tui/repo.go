package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lgtable "github.com/charmbracelet/lipgloss/table"
	"github.com/jjournet/tgr/tui/constants"
)

var (
	symbolMapping = map[string]string{
		"completed":       "✔",
		"pending":         "⚪",
		"failure":         "✖",
		"action_required": "⚠",
		"cancelled":       "🚫",
		"neutral":         "🔘",
		"skipped":         "⏭",
		"stale":           "🕒",
		"success":         "🟢",
		"timed_out":       "⌛",
		"in_progress":     "🟡",
		"queued":          "🟡",
		"requested":       "🟡",
		"waiting":         "🟡",
	}
)

type RepoModel struct {
	top  *lgtable.Table
	main table.Model
}

func InitRepo() (tea.Model, tea.Cmd) {
	m := RepoModel{}
	toprows := [][]string{
		{"Repo", constants.Repo.Name},
		{"Org", constants.Repo.Organization},
	}
	m.top = lgtable.New().Rows(toprows...)
	columns := []table.Column{
		{Title: "Status", Width: 10},
		{Title: "Title", Width: 50},
		{Title: "WorkflowID", Width: 10},
	}
	rows := []table.Row{}
	for _, w := range constants.Repo.GetRuns() {
		log.Printf("Workflow: %v", w.Title)
		rows = append(rows, table.Row{symbolMapping[w.Status], w.Title, fmt.Sprintf("%d", w.WorkflowID)})
	}

	m.main = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	if constants.WindowSize.Height != 0 {
		top, _, bottom, _ := constants.DocStyle.GetMargin()
		m.main.SetHeight(constants.WindowSize.Height - top - bottom - 1 - 5)
	}
	return m, nil
}

func (m RepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, _, bottom, _ := constants.DocStyle.GetMargin()
		m.main.SetHeight(msg.Height - top - bottom - 1 - 5)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "d":
			return InitProfile()
		}
	}
	var cmd tea.Cmd
	m.main, cmd = m.main.Update(msg)
	return m, cmd
}

func (m RepoModel) Init() tea.Cmd {
	return nil
}

func (m RepoModel) View() string {
	var s string
	// style := lipgloss.NewStyle().Width(80).Height(20)
	// get top table height
	s += lipgloss.JoinVertical(lipgloss.Top, m.top.Render(), m.main.View())
	return s
}
