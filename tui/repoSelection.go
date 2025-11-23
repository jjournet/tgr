package tui

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type repoSelection struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// Context
	owner  string
	isUser bool

	// State
	repos   []github.RepoInfo
	loading bool
	err     error

	// UI
	RepoList       table.Model
	visibleCommand bool
}

func (m *repoSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	cmdHeight := 0
	if m.visibleCommand {
		cmdHeight = 3
	}
	slog.Debug("Resizing main", "Header", headerHeight, "Footer", footerHeight, "Command", cmdHeight)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 3 - cmdHeight)
	m.RepoList = m.RepoList.WithPageSize(h - headerHeight - footerHeight - 3 - cmdHeight - 1)
	constants.CommandStyle = constants.CommandStyle.Width(w - 2).Height(1)
}

// NewRepoSelection creates a new repository selection model
func NewRepoSelection(ghService *github.GitHubService, owner string, isUser bool) (tea.Model, tea.Cmd) {
	m := &repoSelection{
		ghService: ghService,
		owner:     owner,
		isUser:    isUser,
		loading:   true,
	}

	m.InitTop("Repository Selection", owner)
	m.TopFields = []string{owner, "Repository Selection", "(Loading...)"}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select", "(/) Filter", "(backspace) Back", "Page: ?"}

	m.visibleCommand = false
	m.CommandInput = textinput.New()

	// Load repos asynchronously
	return m, ghService.LoadReposCmd(owner, isUser)
}

func (m *repoSelection) Init() tea.Cmd {
	return nil
}

func (m *repoSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case github.ReposLoadedMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.repos = msg.Repos
		m.loading = false
		m.TopFields[2] = fmt.Sprintf("(%d repos)", len(m.repos))

		// Build UI table
		m.RepoList = m.buildRepoTable(m.repos)

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

		// Handle filter input mode
		if m.visibleCommand {
			return m.handleFilterInput(msg)
		}

		// Normal key handling
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			repoName := m.RepoList.HighlightedRow().Data["repo"].(string)
			return NewRepoView(m.ghService, m.owner, repoName)
		case "/":
			m.visibleCommand = true
			m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
			m.CommandInput.Focus()
			return m, nil
		case "backspace":
			return NewProfileSelection(m.ghService)
		default:
			var cmd tea.Cmd
			m.RepoList, cmd = m.RepoList.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *repoSelection) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.String() {
	case "enter", "esc":
		m.visibleCommand = false
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
		m.CommandInput.Blur()

		if msg.String() == "esc" {
			m.CommandInput.SetValue("")
			m.RepoList = m.RepoList.WithFilterInputValue("")
			m.TopFields[2] = fmt.Sprintf("(%d repos)", len(m.repos))
		} else {
			m.TopFields[2] = "Filter: " + m.CommandInput.Value()
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.CommandInput, cmd = m.CommandInput.Update(msg)
		cmds = append(cmds, cmd)
		m.RepoList = m.RepoList.WithFilterInput(m.CommandInput)
	}

	return m, tea.Batch(cmds...)
}

func (m *repoSelection) View() string {
	slog.Debug("View Repo Selection")

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit or 'backspace' to go back", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading repositories..."
	}

	// Update arrows
	for i, row := range m.RepoList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.RepoList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}

	m.BottomFields[4] = fmt.Sprintf("Page: %d/%d", m.RepoList.CurrentPage(), m.RepoList.MaxPages())

	if m.visibleCommand {
		return fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			m.RenderTopFields(),
			constants.CommandStyle.BorderForeground(lipgloss.Color("#77c2f9")).Render(m.CommandInput.View()),
			constants.MainStyle.Render(m.RepoList.View()),
			m.RenderBottomFields(),
		)
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.BorderForeground(lipgloss.Color("#77c2f9")).Render(m.RepoList.View()),
		m.RenderBottomFields(),
	)
}

func (m *repoSelection) buildRepoTable(repos []github.RepoInfo) table.Model {
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("repo", "Repository", 40).WithFiltered(true),
		table.NewColumn("desc", "Description", 100),
	}

	rows := []table.Row{}
	for _, repo := range repos {
		rows = append(rows, table.NewRow(table.RowData{
			"repo": repo.Name,
			"desc": repo.Description,
		}))
	}

	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).
		WithFooterVisibility(false).
		WithHighlightedRow(0)
}
