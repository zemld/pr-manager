package db

import (
	"errors"

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

func (s *UserStorage) Select(userID string) (domain.User, error) {
	rows, err := s.tx.Query(s.ctx, s.selectQuery, userID)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	var user domain.User
	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
		if err != nil {
			return domain.User{}, err
		}
	} else {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *UserStorage) Update(user domain.User) error {
	_, err := s.tx.Exec(s.ctx, s.updateQuery, user.IsActive, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) Insert(user domain.User) error {
	_, err := s.tx.Exec(s.ctx, s.insertQuery, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return err
	}
	return nil
}
