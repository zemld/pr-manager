package db

import (
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type TeamStorage struct {
	Config
	Transactor
	selectQuery     string
	insertQuery     string
	selectUserQuery string
	deleteQuery     string
}

func NewTeamStorage(config Config, transactor Transactor) *TeamStorage {
	return &TeamStorage{Config: config, Transactor: transactor}
}

func (s *TeamStorage) SetSelectQuery(selectQuery string) {
	s.selectQuery = selectQuery
}

func (s *TeamStorage) SetInsertQuery(insertQuery string) {
	s.insertQuery = insertQuery
}

func (s *TeamStorage) SetSelectUserQuery(selectUserQuery string) {
	s.selectUserQuery = selectUserQuery
}

func (s *TeamStorage) SetDeleteQuery(deleteQuery string) {
	s.deleteQuery = deleteQuery
}

func (s *TeamStorage) Select(teamName *string) ([]domain.Team, error) {
	var filter any
	if teamName != nil {
		filter = *teamName
	} else {
		filter = nil
	}

	rows, err := s.Transactor.Query(s.ctx, s.selectQuery, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		err = rows.Scan(&team.TeamName, &team.Members)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func (s *TeamStorage) Insert(team domain.Team) error {
	userInserter := NewUserStorage(s.Config, s.Transactor)
	userInserter.SetInsertQuery(s.insertQuery)
	userInserter.SetSelectQuery(SelectUser)

	for _, member := range team.Members {
		userID := member.UserID
		existingUsers, err := userInserter.Select(&userID)
		if err == nil && len(existingUsers) > 0 && existingUsers[0].UserID != "" {
			return domain.ErrUserInAnotherTeam
		}
	}

	for _, member := range team.Members {
		err := userInserter.Insert(domain.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TeamStorage) Delete(teamName string) error {
	_, err := s.Transactor.Exec(s.ctx, s.deleteQuery, teamName)
	if err != nil {
		return err
	}
	return nil
}
