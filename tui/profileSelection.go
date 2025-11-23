package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type profileSelection struct {
	commonElements

	// Service
	ghService *github.GitHubService

	// State
	currentUser string
	orgs        []github.Owner
	owners      []github.Owner
	userLoaded  bool
	orgsLoaded  bool
	loading     bool
	err         error

	// UI
	OwnerList table.Model
}

func (m *profileSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.RenderTopFields())
	footerHeight := lipgloss.Height(m.RenderBottomFields())
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2)
}

// NewProfileSelection creates a new profile selection model
func NewProfileSelection(ghService *github.GitHubService) (tea.Model, tea.Cmd) {
	slog.Debug("NewProfileSelection called")
	m := &profileSelection{
		ghService: ghService,
		loading:   true,
	}

	m.InitTop("Profile Selection", "Loading...")
	m.TopFields = []string{"Profile Selection", "Loading..."}
	m.InitBottom()
	m.BottomFields = []string{"(q) Quit", "(enter) Select"}

	slog.Debug("Returning model with LoadUserCmd and LoadOrgsCmd")
	// Return model and commands to load data
	return m, tea.Batch(
		ghService.LoadUserCmd(),
		ghService.LoadOrgsCmd(),
	)
}

func (m *profileSelection) Init() tea.Cmd {
	return nil
}

func (m *profileSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.UserLoadedMsg:
		slog.Debug("Received UserLoadedMsg", "Login", msg.Login, "Err", msg.Err)
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}
		m.currentUser = msg.Login
		m.userLoaded = true
		m.TopFields = []string{msg.Login, "Profile Selection"}
		m.checkLoadingComplete()
		return m, nil

	case github.OrgsLoadedMsg:
		slog.Debug("Received OrgsLoadedMsg", "Orgs count", len(msg.Orgs), "Err", msg.Err)
		if msg.Err != nil {
			m.err = msg.Err
			m.loading = false
			return m, nil
		}

		m.orgs = msg.Orgs
		m.orgsLoaded = true
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
			return m, nil // Ignore other keys while loading
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedOwner := m.OwnerList.HighlightedRow().Data["profile"].(string)
			isUser := m.OwnerList.HighlightedRow().Data["isUser"].(bool)
			return NewRepoSelection(m.ghService, selectedOwner, isUser)
		}
	}

	if !m.loading {
		var cmd tea.Cmd
		m.OwnerList, cmd = m.OwnerList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *profileSelection) checkLoadingComplete() {
	// Only proceed when both user and orgs are loaded
	if !m.userLoaded || !m.orgsLoaded {
		slog.Debug("Not ready yet", "userLoaded", m.userLoaded, "orgsLoaded", m.orgsLoaded)
		return
	}

	slog.Debug("Both user and orgs loaded, building owner list")

	// Add organizations and current user to owners list
	m.owners = make([]github.Owner, len(m.orgs)+1)
	copy(m.owners, m.orgs)
	m.owners[len(m.orgs)] = github.Owner{
		Login:       m.currentUser,
		Description: "Current User",
		IsUser:      true,
	}

	// Build UI table
	m.OwnerList = m.buildOwnerTable(m.owners)
	m.loading = false

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
}

func (m *profileSelection) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit", m.err)
	}

	if m.loading {
		return m.RenderTopFields() + "\n\nLoading user and organizations..."
	}

	// Update arrows for highlighted row
	for i, row := range m.OwnerList.GetVisibleRows() {
		row.Data["arrow"] = ""
		if i == m.OwnerList.GetHighlightedRowIndex() {
			row.Data["arrow"] = "\uf0a9"
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.RenderTopFields(),
		constants.MainStyle.Render(m.OwnerList.View()),
		m.RenderBottomFields(),
	)
}

func (m *profileSelection) buildOwnerTable(owners []github.Owner) table.Model {
	columns := []table.Column{
		table.NewColumn("arrow", " ", 3),
		table.NewColumn("profile", "Profile", 30),
		table.NewColumn("desc", "Description", 20),
		table.NewColumn("isUser", "", 0), // Hidden column for data
	}

	rows := []table.Row{}
	for _, owner := range owners {
		rows = append(rows, table.NewRow(table.RowData{
			"profile": owner.Login,
			"desc":    owner.Description,
			"isUser":  owner.IsUser,
		}))
	}

	return table.New(columns).WithRows(rows).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true)
}
