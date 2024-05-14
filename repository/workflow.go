package repository

type Workflow struct {
	Name    string
	Repo    string
	Org     string
	State   string
	ID      int64
	HTMLURL string
}

func (w Workflow) GetType() int {
	return WORKFLOW
}

func (w Workflow) GetRepoName() string {
	return w.Repo
}

func (w Workflow) GetOrgName() string {
	return w.Org
}

// func
