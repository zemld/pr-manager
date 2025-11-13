package manager

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

type PullRequestManager struct {
	Storage *db.Storage
}

func NewPullRequestManager(storage *db.Storage) *PullRequestManager {
	return &PullRequestManager{Storage: storage}
}

func (m *PullRequestManager) CreatePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error) {
	err := m.Storage.PullRequestStorage.Create(pullRequest)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return m.Storage.PullRequestStorage.Select(pullRequest.ID)
}

func (m *PullRequestManager) MergePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error) {
	err := m.Storage.PullRequestStorage.Merge(pullRequest)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return m.Storage.PullRequestStorage.Select(pullRequest.ID)
}

func (m *PullRequestManager) ReassignPullRequest(pullRequestID string, oldReviewerID string) (domain.PullRequest, error) {
	pullRequest, err := m.Storage.PullRequestStorage.Select(pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if err := m.isOldReviewerInPullRequest(pullRequest, oldReviewerID); err != nil {
		return domain.PullRequest{}, err
	}

	oldReviewerTeamMembers, err := m.getOldReviewerTeamMembers(oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	newReviewers, err := m.getNewReviewers(oldReviewerTeamMembers, oldReviewerID, pullRequest.AuthorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	err = m.Storage.PullRequestStorage.Reassign(domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID: pullRequestID,
		},
		AssignedReviewers: newReviewers,
	})
	if err != nil {
		return domain.PullRequest{}, err
	}

	return m.Storage.PullRequestStorage.Select(pullRequestID)
}

func (m *PullRequestManager) isOldReviewerInPullRequest(pullRequest domain.PullRequest, oldReviewerID string) error {
	reviewers := strings.Split(strings.Trim(pullRequest.AssignedReviewers, "[]"), ",")
	if !slices.Contains(reviewers, oldReviewerID) {
		return errors.New("old reviewer is not in the list of assigned reviewers")
	}
	return nil
}

func (m *PullRequestManager) getOldReviewerTeamMembers(reviewerID string) ([]domain.TeamMember, error) {
	reviewer, err := m.Storage.UserStorage.Select(reviewerID)
	if err != nil {
		return nil, err
	}

	reviewerTeamName := reviewer.TeamName
	reviewerTeam, err := m.Storage.TeamStorage.Select(reviewerTeamName)
	if err != nil {
		return nil, err
	}

	return reviewerTeam.Members, nil
}

func (m *PullRequestManager) getNewReviewers(teamMembers []domain.TeamMember, oldReviewerID string, authorID string) (string, error) {
	activeUserIDsFromOldReviewerTeamWithoutOldReviewer := make([]string, 0, len(teamMembers))
	for _, member := range teamMembers {
		if member.IsActive && member.UserID != oldReviewerID && member.Username != authorID {
			activeUserIDsFromOldReviewerTeamWithoutOldReviewer = append(
				activeUserIDsFromOldReviewerTeamWithoutOldReviewer,
				member.UserID,
			)
		}
	}

	if len(activeUserIDsFromOldReviewerTeamWithoutOldReviewer) == 0 {
		return "", errors.New("no active users in the old reviewer team")
	}

	newReviewerID := activeUserIDsFromOldReviewerTeamWithoutOldReviewer[rand.IntN(len(activeUserIDsFromOldReviewerTeamWithoutOldReviewer))]
	newReviewers := fmt.Sprintf("[%s, %s]", authorID, newReviewerID)
	return newReviewers, nil
}
