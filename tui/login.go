package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type LoginModel struct {
	authService *github.AuthService
	inputs      []textinput.Model
	focusIndex  int
	err         error
	quitting    bool
	success     bool
}

func NewLoginModel(authService *github.AuthService, username, token string) LoginModel {
	m := LoginModel{
		authService: authService,
		inputs:      make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = constants.CursorStyle
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "GitHub Username"
			t.SetValue(username)
			t.Focus()
			t.PromptStyle = constants.FocusedStyle
			t.TextStyle = constants.FocusedStyle
		case 1:
			t.Placeholder = "GitHub Token"
			t.SetValue(token)
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs)-1 {
				username := m.inputs[0].Value()
				token := m.inputs[1].Value()

				if username == "" || token == "" {
					m.err = fmt.Errorf("username and token are required")
					return m, nil
				}

				err := m.authService.SaveCredentials(username, token)
				if err != nil {
					m.err = err
					return m, nil
				}

				m.success = true
				return m, tea.Quit // We quit this model to signal success to the root
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = constants.FocusedStyle
					m.inputs[i].TextStyle = constants.FocusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = constants.NoStyle
				m.inputs[i].TextStyle = constants.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *LoginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only update the focused input
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	if m.success {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n  Welcome to tgr!\n\n")
	b.WriteString("  Please enter your GitHub credentials.\n")
	b.WriteString("  Token requires 'repo', 'workflow', 'read:org', 'user' scopes.\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &strings.Builder{}
	fmt.Fprintf(button, "\n\n  %s\n\n", constants.FocusedStyle.Render("[ Submit ]"))
	b.WriteString(button.String())

	if m.err != nil {
		b.WriteString(constants.ErrorStyle.Render(fmt.Sprintf("  Error: %v", m.err)))
	}

	return b.String()
}

// GetCredentials returns the entered credentials
func (m LoginModel) GetCredentials() (string, string) {
	return m.inputs[0].Value(), m.inputs[1].Value()
}
