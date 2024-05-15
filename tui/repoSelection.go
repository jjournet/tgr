package tui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

type Repoitem struct {
	Name            string
	RepoDescription string
}

func (i Repoitem) FilterValue() string {
	log.Println("FilterValue: ", i.Name)
	return i.Name
}

func (i Repoitem) Title() string {
	return i.Name
}

func (i Repoitem) Description() string {
	return i.RepoDescription
}

var (
	// itemStyle = lipgloss.NewStyle().PaddingLeft(4)
	// selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.Color("3ceb21"))
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Repoitem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%3d  %-20s %-20s", index+1, i.Name, i.RepoDescription)

	fn := constants.LineStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return constants.SelectedLineStyle.Render("\uf0a9 " + strings.Join(s, " "))
			// return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type repoSelection struct {
	commonElements
	// RepoList table.Model
	RepoList list.Model
}

func (m *repoSelection) resizeMain(w int, h int) {
	headerHeight := lipgloss.Height(m.Top)
	footerHeight := lipgloss.Height(m.Bottom)
	m.RepoList.SetHeight(h - headerHeight - footerHeight - 2 - 1)
	m.RepoList.SetWidth(w - 2)
	constants.MainStyle = constants.MainStyle.Width(w - 2).Height(h - headerHeight - footerHeight - 2 - 5)
	log.Println("new height: ", h, "new width: ", w)
	log.Println("Repo set height: ", h-headerHeight-footerHeight-2-1, "Main set height: ", h-headerHeight-footerHeight-2-5)
}

func InitRepoSelection() (tea.Model, tea.Cmd) {
	m := repoSelection{}
	m.InitTop("Repo Selection", constants.Pr.Profile)
	m.InitBottom()

	var items []list.Item
	for _, repo := range constants.Pr.RepoList {
		items = append(items, list.Item(Repoitem{Name: repo.Name, RepoDescription: repo.Description}))
	}
	// m.RepoList = list.New(items, list.NewDefaultDelegate(), 8, 8)
	m.RepoList = list.New(items, itemDelegate{}, 8, 8)
	m.RepoList.SetShowStatusBar(false)
	m.RepoList.SetShowTitle(false)
	m.RepoList.Styles.PaginationStyle = paginationStyle
	m.RepoList.Styles.HelpStyle = helpStyle

	if constants.WindowSize.Height != 0 {
		m.resizeMain(constants.WindowSize.Width, constants.WindowSize.Height)
	}
	return m, nil
}

func (m repoSelection) Init() tea.Cmd {
	return nil
}

func (m repoSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			repoName := m.RepoList.SelectedItem().(Repoitem).Name
			constants.Repo = repository.NewRepository(repoName, constants.Pr.Profile, constants.User.Client)
			return InitRepo()
		}
	}
	var cmd tea.Cmd
	m.RepoList, cmd = m.RepoList.Update(msg)
	return m, cmd
}

func (m repoSelection) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.Top, constants.MainStyle.Render(m.RepoList.View()), m.Bottom)
}
