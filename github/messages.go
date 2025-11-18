package github

// Messages for the Bubble Tea update cycle
// Each message represents the result of an async operation

// UserLoadedMsg is sent when current user info is loaded
type UserLoadedMsg struct {
	Login string
	Name  string
	Err   error
}

// OrgsLoadedMsg is sent when user's organizations are loaded
type OrgsLoadedMsg struct {
	Orgs []Owner
	Err  error
}

// ReposLoadedMsg is sent when repositories are loaded
type ReposLoadedMsg struct {
	Owner string
	Repos []RepoInfo
	Err   error
}

// RepoDetailsLoadedMsg is sent when detailed repo info is loaded
type RepoDetailsLoadedMsg struct {
	Repo *RepoDetails
	Err  error
}

// WorkflowsLoadedMsg is sent when workflows are loaded
type WorkflowsLoadedMsg struct {
	Workflows []WorkflowInfo
	Err       error
}

// WorkflowRunsLoadedMsg is sent when workflow runs are loaded
type WorkflowRunsLoadedMsg struct {
	WorkflowID int64
	Runs       []RunInfo
	Err        error
}

// RunDetailLoadedMsg is sent when detailed run info is loaded
type RunDetailLoadedMsg struct {
	Run *RunDetailInfo
	Err error
}

// WorkflowTriggeredMsg is sent when a workflow is triggered
type WorkflowTriggeredMsg struct {
	Success bool
	Err     error
}

// IssuesLoadedMsg is sent when issues are loaded
type IssuesLoadedMsg struct {
	Issues []IssueInfo
	Err    error
}

// IssueDetailLoadedMsg is sent when a single issue detail is loaded
type IssueDetailLoadedMsg struct {
	Issue *IssueInfo
	Err   error
}
