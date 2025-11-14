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

const maxAssigners = 2

type PullRequestManager struct {
	Storage *db.Storage
}

func NewPullRequestManager(storage *db.Storage) *PullRequestManager {
	return &PullRequestManager{Storage: storage}
}

func (m *PullRequestManager) CreatePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error) {
	authorTeamMembers, err := m.getReviewerTeamMembers(pullRequest.AuthorID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	possibleAssigners := m.filterReviewers(m.getActiveUserIDsFromTeam(authorTeamMembers), pullRequest.AuthorID)
	if len(possibleAssigners) == 0 {
		return domain.PullRequest{}, errors.New("no possible assigners")
	}
	assigners := make([]string, 0, maxAssigners)
	for range maxAssigners {
		if len(possibleAssigners) == 0 {
			break
		}
		assigners = append(assigners, m.getRandomUserID(possibleAssigners))
		possibleAssigners = m.filterReviewers(possibleAssigners, assigners...)
	}
	pullRequest.AssignedReviewers = fmt.Sprintf("[%s]", strings.Join(assigners, ", "))

	err = m.Storage.PullRequestStorage.Create(pullRequest)
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

	if pullRequest.Status == domain.Merged {
		return domain.PullRequest{}, errors.New("pull request is already merged")
	}

	oldReviewerTeamMembers, err := m.getReviewerTeamMembers(oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	oldReviewers := strings.Split(strings.Trim(pullRequest.AssignedReviewers, "[]"), ",")
	var anotherReviewer string
	for _, reviewer := range oldReviewers {
		if strings.Trim(reviewer, " ") != oldReviewerID {
			anotherReviewer = reviewer
			break
		}
	}

	newPossibleReviewers := m.filterReviewers(m.getActiveUserIDsFromTeam(oldReviewerTeamMembers), oldReviewerID, pullRequest.AuthorID, anotherReviewer)

	newReviewers := make([]string, 0, maxAssigners)
	if anotherReviewer != "" {
		newReviewers = append(newReviewers, anotherReviewer)
	}
	if len(newPossibleReviewers) > 0 {
		newReviewers = append(newReviewers, m.getRandomUserID(newPossibleReviewers))
	}
	pullRequest.AssignedReviewers = fmt.Sprintf("[%s]", strings.Join(newReviewers, ", "))

	err = m.Storage.PullRequestStorage.Reassign(pullRequest)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return m.Storage.PullRequestStorage.Select(pullRequestID)
}

func (m *PullRequestManager) getReviewerTeamMembers(reviewerID string) ([]domain.TeamMember, error) {
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

func (m *PullRequestManager) getActiveUserIDsFromTeam(teamMembers []domain.TeamMember) []string {
	activeUserIDsFromTeam := make([]string, 0, len(teamMembers))
	for _, member := range teamMembers {
		if member.IsActive {
			activeUserIDsFromTeam = append(
				activeUserIDsFromTeam,
				member.UserID,
			)
		}
	}
	return activeUserIDsFromTeam
}

func (m *PullRequestManager) filterReviewers(userIDs []string, destrictedUserIDs ...string) []string {
	filteredUserIDs := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		if !slices.Contains(destrictedUserIDs, userID) {
			filteredUserIDs = append(filteredUserIDs, userID)
		}
	}
	return filteredUserIDs
}

func (m *PullRequestManager) getRandomUserID(userIDs []string) string {
	return userIDs[rand.IntN(len(userIDs))]
}

func (m *PullRequestManager) UserPullRequestsReviews(userID string) ([]domain.PullRequest, error) {
	return m.Storage.PullRequestStorage.SelectUserPullRequestsReviews(userID)
}
