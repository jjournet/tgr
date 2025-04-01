package repositoryviews

import (
	"context"
	"fmt"

	"github.com/evertras/bubble-table/table"
	rep "github.com/jjournet/tgr/repository"
	"github.com/jjournet/tgr/tui/constants"
)

func GetSummaryListModel() table.Model {
	columns := []table.Column{
		table.NewColumn("indicator", " ", 3),
		table.NewColumn("type", "Repository info", 40).WithFiltered(true),
		table.NewColumn("value", "Value", 80),
	}
	repo := constants.Repo.GetRepo()

	var items []table.Row
	// Display Project
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      rep.ConvertRepoElementType(rep.PROJECT),
		"value":     constants.Repo.Name,
		"id":        rep.PROJECT,
	}))
	// Display Description
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      rep.ConvertRepoElementType(rep.DESCRIPTION),
		"value":     fmt.Sprintf("Description: %s", *repo.Description),
		"id":        rep.DESCRIPTION,
	}))
	// Display Workflow
	items = append(items, table.NewRow(table.RowData{
		"indicator": "",
		"type":      rep.ConvertRepoElementType(rep.WORKFLOW),
		"value":     fmt.Sprintf("Workflow: %d", len(constants.Repo.GetWorkflows())),
		"id":        rep.WORKFLOW,
	}))
	items = append(items, table.NewRow(table.RowData{"indicator": "",
		"type":  rep.ConvertRepoElementType(rep.RUN),
		"value": fmt.Sprintf("Actions: %d", len(constants.Repo.Runs)),
		"id":    rep.RUN,
	}))
	// append all languages in one string, with percentage in parenthesis
	var langs string
	// for lang, percent := range constants.Repo.Languages {
	// 	langs += fmt.Sprintf("%s (%d), ", lang, percent)
	// }
	languages, _, _ := constants.Repo.GetClient().Repositories.ListLanguages(context.Background(), constants.Pr.Profile, constants.Repo.Name)
	for lang := range languages {
		langs += fmt.Sprintf("%s (%d) ", lang.Name, lang.Percentage)
	}
	items = append(items, table.NewRow(table.RowData{"indicator": "",
		"type":  "Languages",
		"value": langs,
		"id":    rep.LANGUAGES,
	}))

	return table.New(columns).WithRows(items).
		Focused(true).
		Border(table.Border{}).
		WithBaseStyle(constants.BaseTableStyle).
		HighlightStyle(constants.HighlightedLineStyle).
		Filtered(true).WithHeaderVisibility(false)

}
