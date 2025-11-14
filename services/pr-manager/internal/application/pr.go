package application

import (
	"context"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/manager"
)

func CreatePullRequest(ctx context.Context, pullRequest domain.PullRequest) (domain.PullRequest, error) {
	var result domain.PullRequest
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		pullRequestManager := manager.NewPullRequestManager(configureStorage(tx))
		var err error
		result, err = pullRequestManager.CreatePullRequest(pullRequest)
		return err
	})
	return result, err
}

func MergePullRequest(ctx context.Context, pullRequest domain.PullRequest) (domain.PullRequest, error) {
	var result domain.PullRequest
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		pullRequestManager := manager.NewPullRequestManager(configureStorage(tx))
		var err error
		result, err = pullRequestManager.MergePullRequest(pullRequest)
		return err
	})
	return result, err
}

func ReassignPullRequest(ctx context.Context, pullRequestID string, oldReviewerID string) (domain.PullRequest, error) {
	var result domain.PullRequest
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		pullRequestManager := manager.NewPullRequestManager(configureStorage(tx))
		var err error
		result, err = pullRequestManager.ReassignPullRequest(pullRequestID, oldReviewerID)
		return err
	})
	return result, err
}

func GetUserPullRequestsReviews(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var result []domain.PullRequest
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		pullRequestManager := manager.NewPullRequestManager(configureStorage(tx))
		var err error
		result, err = pullRequestManager.UserPullRequestsReviews(userID)
		return err
	})
	return result, err
}

func configureStorage(tx *db.Transactor) *db.Storage {
	userStorage := db.NewUserStorage(config, *tx)
	userStorage.SetSelectQuery(db.SelectUser)

	teamStorage := db.NewTeamStorage(config, *tx)
	teamStorage.SetSelectQuery(db.SelectTeam)

	pullRequestStorage := db.NewPullRequestStorage(config, *tx)
	pullRequestStorage.SetSelectQuery(db.SelectPullRequest)
	pullRequestStorage.SetCreateQuery(db.CreatePullRequest)
	pullRequestStorage.SetMergeQuery(db.MergePullRequest)
	pullRequestStorage.SetReassignQuery(db.ReassignPullRequest)
	pullRequestStorage.SetUserPullRequestsReviewsQuery(db.UserPullRequestsReviews)

	return &db.Storage{
		Config:             config,
		Transactor:         *tx,
		UserStorage:        userStorage,
		TeamStorage:        teamStorage,
		PullRequestStorage: pullRequestStorage,
	}
}
