package storager

import (
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type Initializer interface {
	Initialize() error
}

type Transactor interface {
	Begin() error
	Commit() error
	Rollback() error
}

type Closer interface {
	Close() error
}

type UserStorager interface {
	UserSelector
	UserUpdater
	UserInserter
}

type UserSelector interface {
	Select(userID string) (domain.User, error)
}

type UserUpdater interface {
	Update(user domain.User) error
}

type UserInserter interface {
	Insert(user domain.User) error
}

type TeamStorager interface {
	TeamSelector
	TeamInserter
}

type TeamSelector interface {
	Select(teamName string) (domain.Team, error)
}

type TeamInserter interface {
	Insert(team domain.Team) error
}

type PullRequestStorager interface {
	PullRequestSelector
	PullRequestUpdater
	PullRequestCreator
	PullRequestMerger
	PullRequestReassigner
}

type PullRequestSelector interface {
	Select(pullRequestID string) (domain.PullRequest, error)
}

type PullRequestUpdater interface {
	Update(pullRequest domain.PullRequest) error
}

type PullRequestCreator interface {
	Create(pullRequest domain.PullRequest) error
}

type PullRequestMerger interface {
	Merge(pullRequest domain.PullRequest) error
}

type PullRequestReassigner interface {
	Reassign(pullRequest domain.PullRequest) error
}
