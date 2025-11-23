package errs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (e ErrorCode) StatusCode() int {
	switch e {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeTeamExists:
		return http.StatusBadRequest
	case CodePRMerged:
		return http.StatusConflict
	case CodePRExists:
		return http.StatusBadRequest
	case CodeNotAssigned:
		return http.StatusConflict
	case CodeNoCandidate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type ApiError struct {
	Code    ErrorCode
	Message string
}

func (e ApiError) Error() string {
	return e.Message
}

func (e ApiError) ReturnError(c *gin.Context, message string, args ...any) {
	var msg string
	if len(args) == 0 {
		msg = message
	} else {
		msg = fmt.Sprintf(message, args...)
	}

	c.JSON(e.Code.StatusCode(), NewErrorResponse(e.Code, msg))
}

func NewApiError(code ErrorCode, message string) ApiError {
	return ApiError{code, message}
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

var ResourceNotFound = NewApiError(CodeNotFound, "resource not found")
var TeamExists = NewApiError(CodeTeamExists, "team exists")
var PullRequestExists = NewApiError(CodePRExists, "pull request exists")
var PullRequestMerged = NewApiError(CodePRMerged, "pull request merged")
var NotAssigned = NewApiError(CodeNotAssigned, "reviewer not assigned")
var NoCandidate = NewApiError(CodeNoCandidate, "no candidate for review")
