package repositoryviews

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	rep "github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

func GetSummaryListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("indicator", " ", 3),
		table.NewColumn("type", "Repository info", 40).WithFiltered(true),
		table.NewColumn("status", "Status", 80),
	}
	items := []table.Row{
		table.NewRow(table.RowData{
			"indicator": "",
			"type":      rep.ConvertRepoElementType(rep.WORKFLOW),
			"status":    fmt.Sprintf("Workflow: %d", len(constants.Repo.Workflows)),
			"id":        rep.WORKFLOW,
		}),
		table.NewRow(table.RowData{
			"indicator": "",
			"type":      rep.ConvertRepoElementType(rep.RUN),
			"status":    fmt.Sprintf("Actions: %d", len(constants.Repo.Runs)),
			"id":        rep.RUN,
		}),
	}
	return table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).WithHeaderVisibility(false)

}
