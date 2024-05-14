package tui

import (
	"fmt"

	"github.com/jjournet/tgr/tui/constants"
)

// func BuildBottom()
type commonElements struct {
	Top    string
	Bottom string
}

func (c *commonElements) InitBottom() {
	c.Bottom = constants.StatusBarStyle.Render("Default Bottom")
}

func (c *commonElements) InitTop(topText ...string) {
	txt := ""
	if len(topText) > 0 {
		for _, t := range topText {
			txt += fmt.Sprintf("> %s ", t)
		}
	} else {
		txt = "Default Top"
	}

	c.Top = constants.TopBarStyle.Render(txt)
}
