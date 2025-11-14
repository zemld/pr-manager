package manager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/storager"
)

type TeamManager struct {
	TeamStorage        storager.TeamStorager
	PullRequestStorage storager.PullRequestStorager
}

func NewTeamManager(teamStorage storager.TeamStorager, pullRequestStorage storager.PullRequestStorager) *TeamManager {
	return &TeamManager{TeamStorage: teamStorage, PullRequestStorage: pullRequestStorage}
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

func (m *TeamManager) DeleteTeam(teamName string) error {
	teamNamePtr := &teamName
	teams, err := m.TeamStorage.Select(teamNamePtr)
	if err != nil {
		return err
	}
	if len(teams) == 0 {
		return errors.New("team not found")
	}
	team := teams[0]

	processedPRs := make(map[string]bool)
	for _, member := range team.Members {
		prs, err := m.PullRequestStorage.SelectUserPullRequestsReviews(member.UserID)
		if err != nil {
			return err
		}

		for _, pr := range prs {
			if processedPRs[pr.ID] {
				continue
			}
			processedPRs[pr.ID] = true

			reviewers := parseReviewers(pr.AssignedReviewers)
			updatedReviewers := make([]string, 0, len(reviewers))
			for _, reviewer := range reviewers {
				isTeamMember := false
				for _, teamMember := range team.Members {
					if reviewer == teamMember.UserID {
						isTeamMember = true
						break
					}
				}
				if !isTeamMember {
					updatedReviewers = append(updatedReviewers, reviewer)
				}
			}

			if len(updatedReviewers) == 0 {
				pr.AssignedReviewers = "[]"
			} else {
				pr.AssignedReviewers = fmt.Sprintf("[%s]", strings.Join(updatedReviewers, ", "))
			}

			err = m.PullRequestStorage.Reassign(pr)
			if err != nil {
				return err
			}
		}
	}

	err = m.TeamStorage.Delete(teamName)
	if err != nil {
		return err
	}
	return nil
}

func parseReviewers(reviewersStr string) []string {
	reviewersStr = strings.Trim(reviewersStr, "[]")
	if reviewersStr == "" {
		return []string{}
	}
	parts := strings.Split(reviewersStr, ",")
	reviewers := make([]string, 0, len(parts))
	for _, part := range parts {
		reviewer := strings.TrimSpace(part)
		if reviewer != "" {
			reviewers = append(reviewers, reviewer)
		}
	}
	return reviewers
}
