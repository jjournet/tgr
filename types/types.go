package types

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
	DESCRIPTION
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
	case DESCRIPTION:
		return "Description"
	default:
		return "Unknown"
	}
}
