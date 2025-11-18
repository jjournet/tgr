package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/github"
	"github.com/jjournet/tgr/tui/constants"
)

type inputField struct {
	name        string
	value       string
	description string
	required    bool
	focused     bool
}

type workflowInputFormView struct {
	// Service
	ghService *github.GitHubService

	// Context
	owner      string
	repoName   string
	workflowID int64

	// State
	branch       string
	branchInput  string
	inputs       []inputField
	focusedIndex int
	triggering   bool
	err          error
	success      bool

	// Return to parent
	parentView tea.Model
}

// NewWorkflowInputForm creates a new workflow input form as an overlay
func NewWorkflowInputForm(ghService *github.GitHubService, owner, repoName string, workflowID int64, parentView tea.Model) tea.Model {
	m := &workflowInputFormView{
		ghService:    ghService,
		owner:        owner,
		repoName:     repoName,
		workflowID:   workflowID,
		branchInput:  "main",
		parentView:   parentView,
		focusedIndex: 0,
	}

	// For now, we'll start with just branch input
	// In a more complete implementation, we'd fetch the workflow file to parse inputs
	return m
}

func (m *workflowInputFormView) Init() tea.Cmd {
	return nil
}

func (m *workflowInputFormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case github.WorkflowTriggeredMsg:
		m.triggering = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.success = true
		return m, nil

	case tea.KeyMsg:
		// If success or error, escape returns to parent
		if m.success || m.err != nil {
			if msg.String() == "esc" || msg.String() == "enter" {
				return m.parentView, nil
			}
			return m, nil
		}

		if m.triggering {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m.parentView, nil

		case "tab", "down":
			m.focusedIndex++
			totalFields := 1 + len(m.inputs) // branch + inputs
			if m.focusedIndex >= totalFields {
				m.focusedIndex = 0
			}

		case "shift+tab", "up":
			m.focusedIndex--
			if m.focusedIndex < 0 {
				totalFields := 1 + len(m.inputs)
				m.focusedIndex = totalFields - 1
			}

		case "enter":
			// Trigger the workflow
			m.triggering = true

			// Build inputs map
			inputsMap := make(map[string]interface{})
			for _, input := range m.inputs {
				if input.value != "" {
					inputsMap[input.name] = input.value
				}
			}

			return m, m.ghService.TriggerWorkflowCmd(m.owner, m.repoName, m.workflowID, m.branchInput, inputsMap)

		case "backspace":
			if m.focusedIndex == 0 {
				if len(m.branchInput) > 0 {
					m.branchInput = m.branchInput[:len(m.branchInput)-1]
				}
			} else {
				idx := m.focusedIndex - 1
				if idx < len(m.inputs) && len(m.inputs[idx].value) > 0 {
					m.inputs[idx].value = m.inputs[idx].value[:len(m.inputs[idx].value)-1]
				}
			}

		default:
			// Add character to focused field
			if len(msg.String()) == 1 {
				if m.focusedIndex == 0 {
					m.branchInput += msg.String()
				} else {
					idx := m.focusedIndex - 1
					if idx < len(m.inputs) {
						m.inputs[idx].value += msg.String()
					}
				}
			}
		}
	}

	return m, nil
}

func (m *workflowInputFormView) View() string {
	// Create the popup
	var popup strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5865F2")).
		Padding(0, 2)

	if m.success {
		popup.WriteString(titleStyle.Render("Workflow Triggered Successfully"))
		popup.WriteString("\n\n")
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22EE82"))
		popup.WriteString(successStyle.Render("✓ Workflow has been queued for execution"))
		popup.WriteString("\n\n")
		popup.WriteString("Press ESC or Enter to continue")
	} else if m.err != nil {
		popup.WriteString(titleStyle.Render("Error Triggering Workflow"))
		popup.WriteString("\n\n")
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
		popup.WriteString(errorStyle.Render(fmt.Sprintf("✗ %v", m.err)))
		popup.WriteString("\n\n")
		popup.WriteString("Press ESC to continue")
	} else if m.triggering {
		popup.WriteString(titleStyle.Render("Trigger Workflow"))
		popup.WriteString("\n\n")
		popup.WriteString("Triggering workflow...")
	} else {
		popup.WriteString(titleStyle.Render("Trigger Workflow"))
		popup.WriteString("\n\n")

		// Branch input
		labelStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B949E")).
			Bold(true)

		focusedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1f2937")).
			Padding(0, 1)

		unfocusedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C9D1D9")).
			Padding(0, 1)

		popup.WriteString(labelStyle.Render("Branch/Ref:"))
		popup.WriteString("\n")
		if m.focusedIndex == 0 {
			popup.WriteString(focusedStyle.Render(m.branchInput + "█"))
		} else {
			popup.WriteString(unfocusedStyle.Render(m.branchInput))
		}
		popup.WriteString("\n\n")

		// Input fields
		for i, input := range m.inputs {
			popup.WriteString(labelStyle.Render(input.name))
			if input.required {
				popup.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("*"))
			}
			popup.WriteString("\n")
			if input.description != "" {
				descStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#6B7280")).
					Italic(true)
				popup.WriteString(descStyle.Render(input.description))
				popup.WriteString("\n")
			}

			if m.focusedIndex == i+1 {
				popup.WriteString(focusedStyle.Render(input.value + "█"))
			} else {
				popup.WriteString(unfocusedStyle.Render(input.value))
			}
			popup.WriteString("\n\n")
		}

		// Instructions
		instrStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
		popup.WriteString(instrStyle.Render("Tab/↓: Next field  Shift+Tab/↑: Previous field"))
		popup.WriteString("\n")
		popup.WriteString(instrStyle.Render("Enter: Trigger  ESC: Cancel"))
	}

	// Style the popup box
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5865F2")).
		Padding(1, 2).
		Width(60)

	popupBox := popupStyle.Render(popup.String())

	// Center the popup
	width := constants.WindowSize.Width
	height := constants.WindowSize.Height

	popupHeight := lipgloss.Height(popupBox)
	popupWidth := lipgloss.Width(popupBox)

	overlayStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	// Position the popup
	verticalPadding := (height - popupHeight) / 2
	horizontalPadding := (width - popupWidth) / 2

	positioned := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popupBox,
		lipgloss.WithWhitespaceChars(" "),
	)

	_ = overlayStyle
	_ = verticalPadding
	_ = horizontalPadding

	// Overlay on parent view
	return positioned
}
