package manager

import (
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

func (m *TeamManager) GetTeam(teamName string) (domain.Team, error) {
	team, err := m.TeamStorage.Select(teamName)
	if err != nil {
		return domain.Team{}, err
	}
	return team, nil
}
