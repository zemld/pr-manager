package db

import (
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type UserStorage struct {
	Config
	Transactor
	selectQuery string
	updateQuery string
	insertQuery string
}

func NewUserStorage(config Config, transactor Transactor) *UserStorage {
	return &UserStorage{Config: config, Transactor: transactor}
}

func (s *UserStorage) SetSelectQuery(selectQuery string) {
	s.selectQuery = selectQuery
}

func (s *UserStorage) SetUpdateQuery(updateQuery string) {
	s.updateQuery = updateQuery
}

func (s *UserStorage) SetInsertQuery(insertQuery string) {
	s.insertQuery = insertQuery
}

func (s *UserStorage) Select(userID *string) ([]domain.User, error) {
	var filter any
	if userID != nil {
		filter = *userID
	} else {
		filter = nil
	}

	rows, err := s.Transactor.Query(s.ctx, s.selectQuery, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) Update(user domain.User) error {
	_, err := s.Transactor.Exec(s.ctx, s.updateQuery, user.IsActive, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) Insert(user domain.User) error {
	_, err := s.Transactor.Exec(s.ctx, s.insertQuery, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return err
	}
	return nil
}
