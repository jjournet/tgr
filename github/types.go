package github

import "time"

// Owner represents a GitHub user or organization
type Owner struct {
	Login       string
	Description string
	IsUser      bool
}

// RepoInfo contains basic repository information
type RepoInfo struct {
	Name        string
	Description string
}

// RepoDetails contains detailed repository information
type RepoDetails struct {
	Name        string
	Description string
	MainBranch  string
	Languages   map[string]int
}

// WorkflowInfo represents a GitHub Actions workflow
type WorkflowInfo struct {
	ID    int64
	Name  string
	State string
	Path  string
}

// RunInfo represents a workflow run
type RunInfo struct {
	ID         int64
	Status     string
	Conclusion string
	Title      string
	Branch     string
	Event      string
	CreatedAt  time.Time
}

// RunDetailInfo represents detailed workflow run information
type RunDetailInfo struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
	Branch     string
	Event      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	RunNumber  int
	RunAttempt int
	HeadSHA    string
	Actor      string
	HTMLURL    string
	JobsURL    string
	LogsURL    string
}

// WorkflowDispatchInputs represents the inputs for a workflow dispatch
type WorkflowDispatchInputs struct {
	Ref    string
	Inputs map[string]interface{}
}

// StepInfo represents a step in a workflow job
type StepInfo struct {
	Name        string
	Status      string
	Conclusion  string
	Number      int
	StartedAt   time.Time
	CompletedAt time.Time
}

// JobInfo represents a job in a workflow run
type JobInfo struct {
	ID          int64
	Name        string
	Status      string
	Conclusion  string
	StartedAt   time.Time
	CompletedAt time.Time
	Steps       []StepInfo
}

// IssueInfo represents a GitHub issue
type IssueInfo struct {
	Number    int
	Title     string
	State     string
	Labels    []string
	Author    string
	Comments  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
}

// WorkflowInputDefinition represents a single input parameter for a workflow
type WorkflowInputDefinition struct {
	Name        string
	Description string
	Required    bool
	Default     string
	Type        string
	Options     []string
}
