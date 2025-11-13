package manager

import (
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type TeamAdder interface {
	AddTeam(team domain.Team) (domain.Team, error)
}

type TeamGetter interface {
	GetTeam(teamName string) (domain.Team, error)
}

type UserUpdater interface {
	UpdateUserStatus(user domain.User) (domain.User, error)
}

type UserSelector interface {
	SelectUser(userID string) (domain.User, error)
}

type PullRequestCreator interface {
	CreatePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error)
}

type PullRequestMerger interface {
	MergePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error)
}

type PullRequestReassigner interface {
	ReassignPullRequest(pullRequestID string, oldReviewerID string) (domain.PullRequest, error)
}

type UserPullRequestReviewer interface {
	UserPullRequestsReviews(userID string) ([]domain.PullRequest, error)
}
