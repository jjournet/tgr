package tui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

type repoSelection struct {
	commonElements
	// RepoList table.Model
	RepoList       table.Model
	visibleCommand bool
}

func (m *repoSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	cmdHeight := 0
	if m.visibleCommand {
		cmdHeight = 3
	}
	log.Printf("Header: %d, Footer: %d, Command: %d\n", headerHeight, footerHeight, cmdHeight)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 3 - cmdHeight)
	m.RepoList = m.RepoList.WithPageSize(h - headerHeight - footerHeight - 3 - cmdHeight - 1)
	constants.CommandStyle = constants.CommandStyle.Width(w - 2).Height(1)
}

func InitRepoSelection() (tea.Model, tea.Cmd) {
	m := repoSelection{}
	m.InitTop("Repo Selection", constants.Pr.Profile)
	m.TopFields = []string{"Repo Selection", constants.Pr.Profile, "(No Filter)"}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select", "(/) Filter", "(backspace) Back", "Page: ?"}
	m.visibleCommand = false
	m.CommandInput = textinput.New()
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("repo", "Repository", 40).WithFiltered(true),
		table.NewColumn("desc", "Description", 100),
	}
	items := []table.Row{}
	for _, repo := range constants.Pr.RepoList {
		items = append(items, table.NewRow(table.RowData{"repo": repo.Name, "desc": repo.Description}))
	}
	m.RepoList = table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).
		WithFooterVisibility(false).
		WithHighlightedRow(0)

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m repoSelection) Init() tea.Cmd {
	return nil
}

func (m repoSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	if m.visibleCommand {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			constants.WindowSize = msg
			m.resizeMain(msg.Width, msg.Height)
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.visibleCommand = false
				m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
				m.CommandInput.Blur()
				cmd = tea.Cmd(func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEsc} })
				cmds = append(cmds, cmd)
			case "esc":
				m.visibleCommand = false
				m.CommandInput.SetValue("")
				m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
				m.CommandInput.Blur()
				m.RepoList = m.RepoList.WithFilterInputValue("")
				m.TopFields[2] = "(No Filter)"
				cmd = tea.Cmd(func() tea.Msg { return tea.KeyMsg{Type: tea.KeyEsc} })
				cmds = append(cmds, cmd)
			default:
				m.CommandInput, cmd = m.CommandInput.Update(msg)
				cmds = append(cmds, cmd)
				m.TopFields[2] = "Filter: " + m.CommandInput.Value()
				m.RepoList = m.RepoList.WithFilterInput(m.CommandInput)
			}
			return m, tea.Batch(cmds...)
		}
	}
	// m.CommandInput, cmdCmd = m.CommandInput.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.resizeMain(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			// get the selected repo
			repoName := m.RepoList.HighlightedRow().Data["repo"].(string)
			constants.Repo = repository.NewRepository(repoName, constants.Pr.Profile, constants.User.Client)
			return InitRepoView()
		case "/":
			m.visibleCommand = true
			m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
			m.CommandInput.Focus()
			return m, nil
		case "backspace":
			return InitProfileSelection()
		default:
			var mainCmd tea.Cmd
			m.RepoList, mainCmd = m.RepoList.Update(msg)
			cmds = append(cmds, mainCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m repoSelection) View() string {
	log.Printf("View Repo Selection\n")
	for i, row := range m.RepoList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.RepoList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}
	m.BottomFields[4] = fmt.Sprintf("Page: %d/%d", m.RepoList.CurrentPage(), m.RepoList.MaxPages())
	if m.visibleCommand {
		return fmt.Sprintf("%s\n%s\n%s\n%s", m.RenderTopFields(), constants.CommandStyle.BorderForeground(lipgloss.Color("#77c2f9")).Render(m.CommandInput.View()), constants.MainStyle.Render(m.RepoList.View()), m.RenderBottomFields())
	} else {
		return fmt.Sprintf("%s\n%s\n%s", m.RenderTopFields(), constants.MainStyle.BorderForeground(lipgloss.Color("#77c2f9")).Render(m.RepoList.View()), m.RenderBottomFields())
	}

}
