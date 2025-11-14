package domain

import "errors"

var (
	ErrNotFound            = errors.New("resource not found")
	ErrPRExists            = errors.New("PR id already exists")
	ErrPRMerged            = errors.New("pull request is already merged")
	ErrTeamExists          = errors.New("team_name already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrReviewerNotFound    = errors.New("reviewer not found")
	ErrNotAssigned         = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrUserInAnotherTeam   = errors.New("user with id is in another team")
	ErrNoPossibleAssigners = errors.New("no possible assigners")
)

type ErrorWithCode struct {
	Err  error
	Code string
}

func (e *ErrorWithCode) Error() string {
	return e.Err.Error()
}

func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}

func NewErrorWithCode(err error, code string) *ErrorWithCode {
	return &ErrorWithCode{Err: err, Code: code}
}
