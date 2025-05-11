package runs

import (
	"context"

	"github.com/google/go-github/v69/github"
)

type Run struct {
	ID         int64
	Status     string
	Title      string
	WorkflowID int64
	Branch     string
	Event      string
	CreatedAt  string
	Conclusion string
}

func NewRun(repoName string, org string, workflow string, id int64, client *github.Client) *Run {
	run, _, err := client.Actions.GetWorkflowRunByID(context.Background(), org, repoName, id)
	if err != nil {
		panic(err)
	}

	return &Run{
		ID:         id,
		Status:     run.GetStatus(),
		Title:      run.GetName(),
		WorkflowID: run.GetWorkflowID(),
		Branch:     run.GetHeadBranch(),
		Event:      run.GetEvent(),
		CreatedAt:  run.GetCreatedAt().String(),
		Conclusion: run.GetConclusion(),
	}
}

func GetRuns(org string, repoName string, workflowID int64, client *github.Client) []*Run {
	runslist, _, err := client.Actions.ListWorkflowRunsByID(context.Background(), org, repoName, workflowID, nil)
	if err != nil {
		panic(err)
	}

	// init constants.Runs
	runs := make([]*Run, 0, len(runslist.WorkflowRuns))

	for _, run := range runslist.WorkflowRuns {
		// log.Printf("Run: %v (id %v)", run.GetName(), run.GetID())
		runs = append(runs, &Run{
			ID:         run.GetID(),
			Status:     run.GetStatus(),
			Conclusion: run.GetConclusion(),
			Title:      run.GetName(),
			WorkflowID: run.GetWorkflowID(),
			Branch:     run.GetHeadBranch(),
			Event:      run.GetEvent(),
			CreatedAt:  run.GetCreatedAt().String(),
		})
	}

	return runs
}
