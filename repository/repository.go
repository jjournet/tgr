package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v69/github"
)

type RepoElement interface {
	GetRepoName() string
	GetOrgName() string
	GetType() int
}

type Workflows struct {
	Name  string
	State string
	ID    int64
}

type Runs struct {
	Status     string
	Title      string
	WorkflowID int64
	Branch     string
	Event      string
	ID         int64
	Date       time.Time
}

type Repository struct {
	Name         string
	Organization string
	client       *github.Client
}

type GHRepoInfo struct {
	RepoName    string
	OrgName     string
	Description string
	HEADCommit  string
	MainBranch  string
	Pending     int
}

func NewRepository(repoName string, org string, client *github.Client) *Repository {
	// repo, _, err := client.Repositories.Get(context.Background(), org, repoName)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Repo: %v (org %v)", repoName, org)

	return &Repository{Name: repoName, Organization: org, client: client}
}

func (r *Repository) GetRepoName() string {
	return r.Name
}

func (r *Repository) GetWorkflows() []Workflows {
	wfs, _, err := r.client.Actions.ListWorkflows(context.Background(), r.Organization, r.Name, nil)
	if err != nil {
		panic(err)
	}
	wfsArr := make([]Workflows, len(wfs.Workflows))
	for i, wf := range wfs.Workflows {
		// log.Printf("Workflow: %v", *wf.Name)
		wfsArr[i] = Workflows{Name: *wf.Name, State: *wf.State, ID: *wf.ID}
	}
	return wfsArr
}

func (r *Repository) GetClient() *github.Client {
	return r.client
}

func (r *Repository) GetRepo() *github.Repository {
	repo, _, err := r.client.Repositories.Get(context.Background(), r.Organization, r.Name)
	if err != nil {
		panic(err)
	}
	return repo
}

func (r *Repository) GetDescription() string {
	repo, _, err := r.client.Repositories.Get(context.Background(), r.Organization, r.Name)
	if err != nil {
		panic(err)
	}
	return *repo.Description
}

func (r *Repository) GetRuns() []Runs {
	runsObj, _, err := r.client.Actions.ListRepositoryWorkflowRuns(context.Background(), r.Organization, r.Name, nil)
	if err != nil {
		panic(err)
	}
	runsArr := make([]Runs, len(runsObj.WorkflowRuns))
	for i, run := range runsObj.WorkflowRuns {
		// log.Printf("Run: %v", *run.Name)
		runsArr[i] = Runs{
			Status:     *run.Status,
			Title:      *run.Name,
			WorkflowID: *run.WorkflowID,
			Branch:     *run.HeadBranch,
			Event:      *run.Event,
			ID:         *run.ID,
			Date:       (*run.CreatedAt).Time,
		}
	}
	return runsArr
}

func (r *Repository) GetLanguages() map[string]int {
	languages, _, err := r.client.Repositories.ListLanguages(context.Background(), r.Organization, r.Name)
	if err != nil {
		panic(err)
	}
	log.Printf("Languages: %v", languages)
	langs := make(map[string]int)
	for lang, size := range languages {
		langs[lang] = size
	}
	return langs
}
