package workflow

import (
	"github.com/google/go-github/v69/github"
	"github.com/jjournet/tgr/types"
)

type Workflow struct {
	Name    string
	Repo    string
	Org     string
	State   string
	ID      int64
	HTMLURL string
	client  *github.Client
}

func (w Workflow) GetType() int {
	return types.WORKFLOW
}

func (w Workflow) GetRepoName() string {
	return w.Repo
}

func (w Workflow) GetOrgName() string {
	return w.Org
}

// func
// parameters (row.Data["value"].(string).constants.Pr.Profile, constants.Repo.GetRepoName(), constants.User.Client
func NewWorkflow(name string, repo string, org string, client *github.Client) *Workflow {
	return &Workflow{Name: name, Repo: repo, Org: org, client: client}
}
