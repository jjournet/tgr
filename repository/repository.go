package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v69/github"
)

const (
	WORKFLOW = iota
	RUN
	ISSUE
	PULL_REQUEST
	BRANCH
	COMMIT
	ENVIRONMENT
	VARIABLE
	PROJECT
	LANGUAGES
)

func ConvertRepoElementType(typeElt int) string {
	switch typeElt {
	case WORKFLOW:
		return "Workflow"
	case RUN:
		return "Run"
	case ISSUE:
		return "Issue"
	case PULL_REQUEST:
		return "Pull Request"
	case BRANCH:
		return "Branch"
	case COMMIT:
		return "Commit"
	case ENVIRONMENT:
		return "Environment"
	case VARIABLE:
		return "Variable"
	case PROJECT:
		return "Project"
	default:
		return "Unknown"
	}
}

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
	Workflows    []Workflows
	Runs         []Runs
	Languages    map[string]int
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
	wfs, _, err := client.Actions.ListWorkflows(context.Background(), org, repoName, nil)
	if err != nil {
		panic(err)
	}
	wfsArr := make([]Workflows, len(wfs.Workflows))
	for i, wf := range wfs.Workflows {
		// log.Printf("Workflow: %v", *wf.Name)
		wfsArr[i] = Workflows{Name: *wf.Name, State: *wf.State, ID: *wf.ID}
	}

	runsObj, _, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), org, repoName, nil)
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
	// retrieve repository languages
	languages, _, err := client.Repositories.ListLanguages(context.Background(), org, repoName)
	if err != nil {
		panic(err)
	}
	log.Printf("Languages: %v", languages)
	langs := make(map[string]int)
	for lang, size := range languages {
		langs[lang] = size
	}
	return &Repository{Name: repoName, Organization: org, Workflows: wfsArr, Runs: runsArr, Languages: langs, client: client}
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
