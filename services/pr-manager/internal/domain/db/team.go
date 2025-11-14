package db

import (
	"fmt"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type TeamStorage struct {
	Config
	Transactor
	selectQuery     string
	insertQuery     string
	selectUserQuery string
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

func (s *TeamStorage) Select(teamName string) (domain.Team, error) {
	rows, err := s.Transactor.Query(s.ctx, s.selectQuery, teamName)
	if err != nil {
		return domain.Team{}, err
	}
	defer rows.Close()

	var team domain.Team
	if rows.Next() {
		err = rows.Scan(&team.TeamName, &team.Members)
		if err != nil {
			return domain.Team{}, err
		}
	}
	return team, nil
}

func (s *TeamStorage) Insert(team domain.Team) error {
	userInserter := NewUserStorage(s.Config, s.Transactor)
	userInserter.SetInsertQuery(s.insertQuery)
	userInserter.SetSelectQuery(SelectUser)

	for _, member := range team.Members {
		existingUser, err := userInserter.Select(member.UserID)
		if err == nil && existingUser.UserID != "" {
			return fmt.Errorf("user with id %s is in another team", member.UserID)
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
