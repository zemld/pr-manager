package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Initializer struct {
	Config  Config
	Queries []string
}

func NewDBInitializer(config Config, queries []string) *Initializer {
	return &Initializer{Config: config, Queries: queries}
}

func (d *Initializer) Initialize() error {
	conn, err := pgx.Connect(context.Background(), d.Config.GetConnectionString())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for _, query := range d.Queries {
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			return err
		}
	}
	return nil
}
