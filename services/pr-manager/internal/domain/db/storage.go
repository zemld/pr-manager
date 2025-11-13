package db

type Storage struct {
	Config
	Transactor
	UserStorage
	TeamStorage
	PullRequestStorage
}

func NewStorage(config Config, transactor Transactor) *Storage {
	return &Storage{
		Config:             config,
		Transactor:         transactor,
		UserStorage:        *NewUserStorage(config, transactor),
		TeamStorage:        *NewTeamStorage(config, transactor),
		PullRequestStorage: *NewPullRequestStorage(config, transactor),
	}
}
