package runs

type Runs struct {
	ID         int64
	Status     string
	Title      string
	WorkflowID int64
	Branch     string
	Event      string
	CreatedAt  string
}

func NewRun(repoName string, org string, workflow string. client *github.Client) *Runs {
	return &Runs{
		ID:         0,
		Status:     "",
		Title:      "",
		WorkflowID: 0,
		Branch:     "",
		Event:      "",
		CreatedAt:  "",
	}
}