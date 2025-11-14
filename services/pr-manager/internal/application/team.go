package application

import (
	"context"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/manager"
)

func AddTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	var result domain.Team
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		teamStorage := db.NewTeamStorage(config, *tx)
		teamStorage.SetInsertQuery(db.InsertUser)
		teamStorage.SetSelectUserQuery(db.SelectUser)
		teamManager := manager.NewTeamManager(teamStorage)
		var err error
		result, err = teamManager.AddTeam(team)
		return err
	}, false)
	return result, err
}

func GetTeam(ctx context.Context, teamName *string) (domain.Team, error) {
	var result domain.Team
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		teamStorage := db.NewTeamStorage(config, *tx)
		teamStorage.SetSelectQuery(db.SelectTeam)
		teamManager := manager.NewTeamManager(teamStorage)
		var err error
		result, err = teamManager.GetTeam(teamName)
		return err
	}, true)
	return result, err
}

func GetTeams(ctx context.Context) ([]domain.Team, error) {
	var result []domain.Team
	err := executor.withTransaction(ctx, func(tx *db.Transactor) error {
		teamStorage := db.NewTeamStorage(config, *tx)
		teamStorage.SetSelectQuery(db.SelectTeam)
		teamManager := manager.NewTeamManager(teamStorage)
		var err error
		result, err = teamManager.GetTeams(nil)
		return err
	}, true)
	return result, err
}
