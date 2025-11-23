package errs

import "errors"

type ErrorCode int

const (
	CodeNotFound = iota
	CodeTeamExists
	CodePRExists
	CodePRMerged
	CodeNotAssigned
	CodeNoCandidate
)

func (e ErrorCode) String() string {
	switch e {
	case CodeNotFound:
		return "NOT_FOUND"
	case CodeTeamExists:
		return "TEAM_EXISTS"
	case CodePRMerged:
		return "PR_MERGED"
	case CodePRExists:
		return "PR_EXISTS"
	case CodeNotAssigned:
		return "NOT_ASSIGNED"
	case CodeNoCandidate:
		return "NO_CANDIDATE"
	default:
		return "ERROR"
	}
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorBody{
			Code:    code.String(),
			Message: message,
		},
	}
}

var ResourceNotFound = errors.New("resource not found")
var TeamExists = errors.New("team exists")
var PullRequestExists = errors.New("pull request exists")
var PullRequestMerged = errors.New("pull request merged")
var NotAssigned = errors.New("reviewer not assigned")
var NoCandidate = errors.New("no candidate for review")
