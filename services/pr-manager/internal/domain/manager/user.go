package manager

import (
	"errors"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/storager"
)

type UserManager struct {
	UserStorage storager.UserStorager
}

func NewUserManager(userStorage storager.UserStorager) *UserManager {
	return &UserManager{UserStorage: userStorage}
}

func (m *UserManager) UpdateUserStatus(user domain.User) (domain.User, error) {
	existingUser, err := m.SelectUser(&user.UserID)
	if err != nil {
		return domain.User{}, err
	}

	existingUser.IsActive = user.IsActive

	err = m.UserStorage.Update(existingUser)
	if err != nil {
		return domain.User{}, err
	}

	return existingUser, nil
}

func (m *UserManager) SelectUser(userID *string) (domain.User, error) {
	users, err := m.UserStorage.Select(userID)
	if err != nil {
		return domain.User{}, err
	}
	if len(users) == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return users[0], nil
}

func (m *UserManager) SelectUsers(userID *string) ([]domain.User, error) {
	return m.UserStorage.Select(userID)
}
