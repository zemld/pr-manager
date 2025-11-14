package db

import "github.com/zemld/pr-manager/pr-manager/internal/domain"

type TeamStorage struct {
	Config
	Transactor
	selectQuery string
	insertQuery string
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
