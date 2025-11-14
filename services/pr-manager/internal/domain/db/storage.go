package db

import "github.com/zemld/pr-manager/pr-manager/internal/domain/storager"

type Storage struct {
	Config
	Transactor
	UserStorage        storager.UserStorager
	TeamStorage        storager.TeamStorager
	PullRequestStorage storager.PullRequestStorager
}

func NewStorage(config Config, transactor Transactor) *Storage {
	return &Storage{
		Config:             config,
		Transactor:         transactor,
		UserStorage:        NewUserStorage(config, transactor),
		TeamStorage:        NewTeamStorage(config, transactor),
		PullRequestStorage: NewPullRequestStorage(config, transactor),
	}
}
