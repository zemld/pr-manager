package manager

import (
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
	user, err := m.SelectUser(user.UserID)
	if err != nil {
		return domain.User{}, err
	}

	err = m.UserStorage.Update(user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (m *UserManager) SelectUser(userID string) (domain.User, error) {
	user, err := m.UserStorage.Select(userID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
