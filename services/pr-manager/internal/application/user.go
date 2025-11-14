package application

import (
	"context"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/manager"
)

func UpdateUserStatus(ctx context.Context, user domain.User) (domain.User, error) {
	var updatedUser domain.User
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		userStorage := db.NewUserStorage(config, *tx)
		userStorage.SetUpdateQuery(db.UpdateUserStatus)
		userStorage.SetSelectQuery(db.SelectUser)
		userManager := manager.NewUserManager(userStorage)
		var err error
		updatedUser, err = userManager.UpdateUserStatus(user)
		return err
	})
	return updatedUser, err
}
