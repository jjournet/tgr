package common

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/jjournet/tgr/tui/constants"
)

// func BuildBottom()
type commonElements struct {
	Top          string
	TopFields    []string
	Bottom       string
	BottomFields []string
	CommandInput textinput.Model
}

var statusStyle = lipgloss.NewStyle().
	Inherit(constants.StatusBarStyle).
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#FF5F87")).
	Padding(0, 1).
	MarginRight(1)

func (c *commonElements) InitBottom() {
	// cat \ueeed , github \uf09b \uf113
	c.Bottom = statusStyle.Render("\uf113") + constants.StatusBarStyle.Render("Default Bottom")
}

func (c *commonElements) InitTop(topText ...string) {
	txt := ""
	if len(topText) > 0 {
		for _, t := range topText {
			txt += fmt.Sprintf(" %s ", t)
		}
	} else {
		txt = "Default Top"
	}

	c.Top = constants.TopBarStyle.Render(txt)
}

func (c *commonElements) RenderTopFields() string {
	var aggregated string
	for i := 0; i < len(c.TopFields); i++ {
		aggregated += " " + c.TopFields[i] + " "
	}
	return constants.TopBarStyle.Render(aggregated)
}

func (c *commonElements) RenderBottomFields() string {
	var aggregated string
	for i := 0; i < len(c.BottomFields); i++ {
		aggregated += " " + c.BottomFields[i] + " "
	}
	return constants.StatusBarStyle.Render(aggregated)
}
