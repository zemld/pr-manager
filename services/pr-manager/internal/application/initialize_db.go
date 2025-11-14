package application

import (
	"context"

	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

func InitializeDB(ctx context.Context) error {
	return db.NewDBInitializer(config,
		db.CreateUsersTable,
		db.CreatePullRequestsStatusesTable,
		db.CreatePullRequestsTable,
		db.FillPullRequestsStatusesTable,
	).Initialize()
}
