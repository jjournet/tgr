package github

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/google/go-github/v69/github"
	"gopkg.in/yaml.v3"
)

// LoadUserCmd returns a command that loads the current user's information
func (s *GitHubService) LoadUserCmd() tea.Cmd {
	return func() tea.Msg {
		log.Println("LoadUserCmd: Starting to fetch user info...")
		user, _, err := s.client.Users.Get(s.Context(), "")
		if err != nil {
			log.Printf("LoadUserCmd: Error fetching user: %v", err)
			return UserLoadedMsg{Err: err}
		}

		log.Printf("LoadUserCmd: Successfully loaded user: %s", user.GetLogin())
		return UserLoadedMsg{
			Login: user.GetLogin(),
			Name:  user.GetName(),
			Err:   nil,
		}
	}
}

// LoadOrgsCmd returns a command that loads the user's organizations
func (s *GitHubService) LoadOrgsCmd() tea.Cmd {
	return func() tea.Msg {
		log.Println("LoadOrgsCmd: Starting to fetch organizations...")
		orgs, _, err := s.client.Organizations.List(s.Context(), "", nil)
		if err != nil {
			log.Printf("LoadOrgsCmd: Error fetching orgs: %v", err)
			return OrgsLoadedMsg{Err: err}
		}

		log.Printf("LoadOrgsCmd: Successfully loaded %d organizations", len(orgs))
		owners := make([]Owner, len(orgs))
		for i, org := range orgs {
			desc := org.GetDescription()
			if desc == "" {
				desc = "Organization"
			}
			owners[i] = Owner{
				Login:       org.GetLogin(),
				Description: desc,
				IsUser:      false,
			}
		}

		return OrgsLoadedMsg{
			Orgs: owners,
			Err:  nil,
		}
	}
}

// LoadReposCmd returns a command that loads repositories for an owner
func (s *GitHubService) LoadReposCmd(owner string, isUser bool) tea.Cmd {
	return func() tea.Msg {
		var allRepos []*gh.Repository
		listOpt := &gh.ListOptions{PerPage: 100}

		for {
			var repos []*gh.Repository
			var resp *gh.Response
			var err error

			if isUser {
				opts := &gh.RepositoryListByUserOptions{ListOptions: *listOpt}
				repos, resp, err = s.client.Repositories.ListByUser(
					s.Context(),
					owner,
					opts,
				)
			} else {
				opts := &gh.RepositoryListByOrgOptions{ListOptions: *listOpt}
				repos, resp, err = s.client.Repositories.ListByOrg(
					s.Context(),
					owner,
					opts,
				)
			}

			if err != nil {
				return ReposLoadedMsg{Owner: owner, Err: err}
			}

			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}
			listOpt.Page = resp.NextPage
		}

		repoInfos := make([]RepoInfo, len(allRepos))
		for i, repo := range allRepos {
			desc := repo.GetDescription()
			if desc == "" {
				desc = "N/A"
			}
			repoInfos[i] = RepoInfo{
				Name:        repo.GetName(),
				Description: desc,
			}
		}

		return ReposLoadedMsg{
			Owner: owner,
			Repos: repoInfos,
			Err:   nil,
		}
	}
}

// LoadRepoDetailsCmd returns a command that loads detailed repo information
func (s *GitHubService) LoadRepoDetailsCmd(owner, repoName string) tea.Cmd {
	return func() tea.Msg {
		repo, _, err := s.client.Repositories.Get(s.Context(), owner, repoName)
		if err != nil {
			return RepoDetailsLoadedMsg{Err: err}
		}

		languages, _, err := s.client.Repositories.ListLanguages(
			s.Context(),
			owner,
			repoName,
		)
		if err != nil {
			languages = make(map[string]int)
		}

		return RepoDetailsLoadedMsg{
			Repo: &RepoDetails{
				Name:        repo.GetName(),
				Description: repo.GetDescription(),
				MainBranch:  repo.GetDefaultBranch(),
				Languages:   languages,
			},
			Err: nil,
		}
	}
}

// LoadWorkflowsCmd returns a command that loads workflows for a repository
func (s *GitHubService) LoadWorkflowsCmd(owner, repoName string) tea.Cmd {
	return func() tea.Msg {
		workflows, _, err := s.client.Actions.ListWorkflows(
			s.Context(),
			owner,
			repoName,
			nil,
		)
		if err != nil {
			return WorkflowsLoadedMsg{Err: err}
		}

		infos := make([]WorkflowInfo, len(workflows.Workflows))
		for i, wf := range workflows.Workflows {
			infos[i] = WorkflowInfo{
				ID:    wf.GetID(),
				Name:  wf.GetName(),
				State: wf.GetState(),
				Path:  wf.GetPath(),
			}
		}

		return WorkflowsLoadedMsg{
			Workflows: infos,
			Err:       nil,
		}
	}
}

// LoadWorkflowRunsCmd returns a command that loads runs for a specific workflow
func (s *GitHubService) LoadWorkflowRunsCmd(owner, repoName string, workflowID int64) tea.Cmd {
	return func() tea.Msg {
		runs, _, err := s.client.Actions.ListWorkflowRunsByID(
			s.Context(),
			owner,
			repoName,
			workflowID,
			nil,
		)
		if err != nil {
			return WorkflowRunsLoadedMsg{WorkflowID: workflowID, Err: err}
		}

		infos := make([]RunInfo, len(runs.WorkflowRuns))
		for i, run := range runs.WorkflowRuns {
			infos[i] = RunInfo{
				ID:         run.GetID(),
				Status:     run.GetStatus(),
				Conclusion: run.GetConclusion(),
				Title:      run.GetName(),
				Branch:     run.GetHeadBranch(),
				Event:      run.GetEvent(),
				CreatedAt:  run.GetCreatedAt().Time,
			}
		}

		return WorkflowRunsLoadedMsg{
			WorkflowID: workflowID,
			Runs:       infos,
			Err:        nil,
		}
	}
}

// LoadAllRepoRunsCmd returns a command that loads all workflow runs for a repo
func (s *GitHubService) LoadAllRepoRunsCmd(owner, repoName string) tea.Cmd {
	return func() tea.Msg {
		runs, _, err := s.client.Actions.ListRepositoryWorkflowRuns(
			s.Context(),
			owner,
			repoName,
			nil,
		)
		if err != nil {
			return WorkflowRunsLoadedMsg{Err: err}
		}

		infos := make([]RunInfo, len(runs.WorkflowRuns))
		for i, run := range runs.WorkflowRuns {
			infos[i] = RunInfo{
				ID:         run.GetID(),
				Status:     run.GetStatus(),
				Conclusion: run.GetConclusion(),
				Title:      run.GetName(),
				Branch:     run.GetHeadBranch(),
				Event:      run.GetEvent(),
				CreatedAt:  run.GetCreatedAt().Time,
			}
		}

		return WorkflowRunsLoadedMsg{
			Runs: infos,
			Err:  nil,
		}
	}
}

// LoadIssuesCmd returns a command that loads issues for a repository
func (s *GitHubService) LoadIssuesCmd(owner, repoName string) tea.Cmd {
	return func() tea.Msg {
		log.Println("[LoadIssuesCmd] Starting to load issues for", owner, "/", repoName)

		issues, _, err := s.client.Issues.ListByRepo(
			s.Context(),
			owner,
			repoName,
			&gh.IssueListByRepoOptions{
				State: "all",
			},
		)
		if err != nil {
			log.Println("[LoadIssuesCmd] Error loading issues:", err)
			return IssuesLoadedMsg{Err: err}
		}

		log.Println("[LoadIssuesCmd] Successfully loaded", len(issues), "issues")

		infos := make([]IssueInfo, len(issues))
		for i, issue := range issues {
			// Skip pull requests (they appear in issue list but have PullRequestLinks)
			if issue.IsPullRequest() {
				continue
			}

			labels := make([]string, len(issue.Labels))
			for j, label := range issue.Labels {
				labels[j] = label.GetName()
			}

			author := ""
			if issue.User != nil {
				author = issue.User.GetLogin()
			}

			infos[i] = IssueInfo{
				Number:    issue.GetNumber(),
				Title:     issue.GetTitle(),
				State:     issue.GetState(),
				Labels:    labels,
				Author:    author,
				Comments:  issue.GetComments(),
				CreatedAt: issue.GetCreatedAt().Time,
				UpdatedAt: issue.GetUpdatedAt().Time,
				Body:      issue.GetBody(),
			}
		}

		return IssuesLoadedMsg{
			Issues: infos,
			Err:    nil,
		}
	}
}

// LoadRunDetailCmd returns a command that loads detailed information for a workflow run
func (s *GitHubService) LoadRunDetailCmd(owner, repoName string, runID int64) tea.Cmd {
	return func() tea.Msg {
		log.Println("[LoadRunDetailCmd] Starting to load run detail for run ID:", runID)

		run, _, err := s.client.Actions.GetWorkflowRunByID(
			s.Context(),
			owner,
			repoName,
			runID,
		)
		if err != nil {
			log.Println("[LoadRunDetailCmd] Error loading run detail:", err)
			return RunDetailLoadedMsg{Err: err}
		}

		log.Println("[LoadRunDetailCmd] Successfully loaded run detail")

		actor := ""
		if run.Actor != nil {
			actor = run.Actor.GetLogin()
		}

		detail := &RunDetailInfo{
			ID:         run.GetID(),
			Name:       run.GetName(),
			Status:     run.GetStatus(),
			Conclusion: run.GetConclusion(),
			Branch:     run.GetHeadBranch(),
			Event:      run.GetEvent(),
			CreatedAt:  run.GetCreatedAt().Time,
			UpdatedAt:  run.GetUpdatedAt().Time,
			RunNumber:  run.GetRunNumber(),
			RunAttempt: run.GetRunAttempt(),
			HeadSHA:    run.GetHeadSHA(),
			Actor:      actor,
			HTMLURL:    run.GetHTMLURL(),
			JobsURL:    run.GetJobsURL(),
			LogsURL:    run.GetLogsURL(),
		}

		return RunDetailLoadedMsg{
			Run: detail,
			Err: nil,
		}
	}
}

// LoadRunJobsCmd returns a command that loads jobs for a workflow run
func (s *GitHubService) LoadRunJobsCmd(owner, repoName string, runID int64) tea.Cmd {
	return func() tea.Msg {
		log.Println("[LoadRunJobsCmd] Starting to load jobs for run ID:", runID)

		jobs, _, err := s.client.Actions.ListWorkflowJobs(
			s.Context(),
			owner,
			repoName,
			runID,
			nil,
		)
		if err != nil {
			log.Println("[LoadRunJobsCmd] Error loading jobs:", err)
			return RunJobsLoadedMsg{Err: err}
		}

		log.Println("[LoadRunJobsCmd] Successfully loaded", len(jobs.Jobs), "jobs")

		jobInfos := make([]JobInfo, len(jobs.Jobs))
		for i, job := range jobs.Jobs {
			steps := make([]StepInfo, len(job.Steps))
			for j, step := range job.Steps {
				steps[j] = StepInfo{
					Name:        step.GetName(),
					Status:      step.GetStatus(),
					Conclusion:  step.GetConclusion(),
					Number:      int(step.GetNumber()),
					StartedAt:   step.GetStartedAt().Time,
					CompletedAt: step.GetCompletedAt().Time,
				}
			}

			jobInfos[i] = JobInfo{
				ID:          job.GetID(),
				Name:        job.GetName(),
				Status:      job.GetStatus(),
				Conclusion:  job.GetConclusion(),
				StartedAt:   job.GetStartedAt().Time,
				CompletedAt: job.GetCompletedAt().Time,
				Steps:       steps,
			}
		}

		return RunJobsLoadedMsg{
			RunID: runID,
			Jobs:  jobInfos,
			Err:   nil,
		}
	}
}

// TriggerWorkflowCmd returns a command that triggers a workflow dispatch event
func (s *GitHubService) TriggerWorkflowCmd(owner, repoName string, workflowID int64, ref string, inputs map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		log.Printf("[TriggerWorkflowCmd] Triggering workflow %d on ref %s", workflowID, ref)

		event := gh.CreateWorkflowDispatchEventRequest{
			Ref:    ref,
			Inputs: inputs,
		}

		_, err := s.client.Actions.CreateWorkflowDispatchEventByID(
			s.Context(),
			owner,
			repoName,
			workflowID,
			event,
		)

		if err != nil {
			log.Println("[TriggerWorkflowCmd] Error triggering workflow:", err)
			return WorkflowTriggeredMsg{Success: false, Err: err}
		}

		log.Println("[TriggerWorkflowCmd] Successfully triggered workflow")
		return WorkflowTriggeredMsg{Success: true, Err: nil}
	}
}

// LoadWorkflowInputsCmd loads the inputs for a workflow
func (s *GitHubService) LoadWorkflowInputsCmd(owner, repoName string, workflowPath string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("[LoadWorkflowInputsCmd] Loading inputs for %s", workflowPath)

		// Get file content
		fileContent, _, _, err := s.client.Repositories.GetContents(
			s.Context(),
			owner,
			repoName,
			workflowPath,
			nil,
		)
		if err != nil {
			log.Printf("[LoadWorkflowInputsCmd] Error fetching file content: %v", err)
			return WorkflowInputsLoadedMsg{Err: err}
		}

		content, err := fileContent.GetContent()
		if err != nil {
			log.Printf("[LoadWorkflowInputsCmd] Error decoding file content: %v", err)
			return WorkflowInputsLoadedMsg{Err: err}
		}

		// Parse YAML
		var wf struct {
			On struct {
				WorkflowDispatch struct {
					Inputs map[string]struct {
						Description string      `yaml:"description"`
						Required    bool        `yaml:"required"`
						Default     interface{} `yaml:"default"` // Default can be string or boolean
						Type        string      `yaml:"type"`
						Options     []string    `yaml:"options"`
					} `yaml:"inputs"`
				} `yaml:"workflow_dispatch"`
			} `yaml:"on"`
		}

		if err := yaml.Unmarshal([]byte(content), &wf); err != nil {
			// If unmarshal fails, it might be because 'on' is not a map.
			// In that case, there are no inputs for workflow_dispatch (or it's not enabled).
			log.Printf("[LoadWorkflowInputsCmd] YAML unmarshal failed (likely no inputs): %v", err)
			return WorkflowInputsLoadedMsg{Inputs: []WorkflowInputDefinition{}, Err: nil}
		}

		var inputs []WorkflowInputDefinition
		for name, input := range wf.On.WorkflowDispatch.Inputs {
			defVal := ""
			if input.Default != nil {
				defVal = fmt.Sprintf("%v", input.Default)
			}

			inputs = append(inputs, WorkflowInputDefinition{
				Name:        name,
				Description: input.Description,
				Required:    input.Required,
				Default:     defVal,
				Type:        input.Type,
				Options:     input.Options,
			})
		}

		return WorkflowInputsLoadedMsg{Inputs: inputs, Err: nil}
	}
}

// FindLatestRunCmd finds the latest run for a workflow
func (s *GitHubService) FindLatestRunCmd(owner, repoName string, workflowID int64) tea.Cmd {
	return func() tea.Msg {
		// Wait a bit to allow GitHub to create the run
		time.Sleep(2 * time.Second)

		runs, _, err := s.client.Actions.ListWorkflowRunsByID(
			s.Context(),
			owner,
			repoName,
			workflowID,
			&gh.ListWorkflowRunsOptions{ListOptions: gh.ListOptions{PerPage: 1}},
		)
		if err != nil {
			return LatestRunFoundMsg{Err: err}
		}

		if len(runs.WorkflowRuns) == 0 {
			return LatestRunFoundMsg{Err: fmt.Errorf("no runs found")}
		}

		run := runs.WorkflowRuns[0]
		return LatestRunFoundMsg{
			RunID: run.GetID(),
			Err:   nil,
		}
	}
}
