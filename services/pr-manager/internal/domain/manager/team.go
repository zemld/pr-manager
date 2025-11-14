package manager

import (
	"errors"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/storager"
)

type TeamManager struct {
	TeamStorage storager.TeamStorager
}

func NewTeamManager(teamStorage storager.TeamStorager) *TeamManager {
	return &TeamManager{TeamStorage: teamStorage}
}

func (m *TeamManager) AddTeam(team domain.Team) (domain.Team, error) {
	err := m.TeamStorage.Insert(team)
	if err != nil {
		return domain.Team{}, err
	}
	return team, nil
}

func (m *TeamManager) GetTeam(teamName *string) (domain.Team, error) {
	teams, err := m.TeamStorage.Select(teamName)
	if err != nil {
		return domain.Team{}, err
	}
	if len(teams) == 0 {
		return domain.Team{}, errors.New("team not found")
	}
	return teams[0], nil
}

func (m *TeamManager) GetTeams(teamName *string) ([]domain.Team, error) {
	return m.TeamStorage.Select(teamName)
}
