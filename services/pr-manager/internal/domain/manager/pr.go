package manager

import (
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
	assigners := make([]string, 0, maxAssigners)
	for range maxAssigners {
		if len(possibleAssigners) == 0 {
			break
		}
		assigners = append(assigners, m.getRandomUserID(possibleAssigners))
		possibleAssigners = m.filterReviewers(possibleAssigners, assigners...)
	}
	if len(assigners) == 0 {
		pullRequest.AssignedReviewers = "[]"
	} else {
		pullRequest.AssignedReviewers = fmt.Sprintf("[%s]", strings.Join(assigners, ", "))
	}

	err = m.Storage.PullRequestStorage.Create(pullRequest)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return domain.PullRequest{}, domain.ErrPRExists
		}
		return domain.PullRequest{}, err
	}

	prs, err := m.Storage.PullRequestStorage.Select(&pullRequest.ID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	if len(prs) == 0 {
		return domain.PullRequest{}, domain.ErrNotFound
	}
	return prs[0], nil
}

func (m *PullRequestManager) MergePullRequest(pullRequest domain.PullRequest) (domain.PullRequest, error) {
	err := m.Storage.PullRequestStorage.Merge(pullRequest)
	if err != nil {
		return domain.PullRequest{}, err
	}

	prs, err := m.Storage.PullRequestStorage.Select(&pullRequest.ID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	if len(prs) == 0 {
		return domain.PullRequest{}, domain.ErrNotFound
	}
	return prs[0], nil
}

func (m *PullRequestManager) ReassignPullRequest(pullRequestID string, oldReviewerID string) (domain.PullRequest, string, error) {
	prs, err := m.Storage.PullRequestStorage.Select(&pullRequestID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}
	if len(prs) == 0 {
		return domain.PullRequest{}, "", domain.ErrNotFound
	}
	pullRequest := prs[0]

	if pullRequest.Status == domain.Merged {
		return domain.PullRequest{}, "", domain.ErrPRMerged
	}

	oldReviewerTeamMembers, err := m.getReviewerTeamMembers(oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	oldReviewers := strings.Split(strings.Trim(pullRequest.AssignedReviewers, "[]"), ",")
	var anotherReviewer string
	foundOldReviewer := false
	for _, reviewer := range oldReviewers {
		reviewer = strings.Trim(reviewer, " ")
		if reviewer == oldReviewerID {
			foundOldReviewer = true
		} else if reviewer != "" {
			anotherReviewer = reviewer
		}
	}
	if !foundOldReviewer {
		return domain.PullRequest{}, "", domain.ErrNotAssigned
	}

	newPossibleReviewers := m.filterReviewers(m.getActiveUserIDsFromTeam(oldReviewerTeamMembers), oldReviewerID, pullRequest.AuthorID, anotherReviewer)

	updatedReviewers := make([]string, 0, maxAssigners)
	if anotherReviewer != "" {
		updatedReviewers = append(updatedReviewers, anotherReviewer)
	}
	newReviewer := ""
	if len(newPossibleReviewers) > 0 {
		newReviewer = m.getRandomUserID(newPossibleReviewers)
		updatedReviewers = append(updatedReviewers, newReviewer)
	} else if len(updatedReviewers) == 0 {
		// No replacement candidate available and no other reviewer
		return domain.PullRequest{}, "", domain.ErrNoCandidate
	}
	pullRequest.AssignedReviewers = fmt.Sprintf("[%s]", strings.Join(updatedReviewers, ", "))

	err = m.Storage.PullRequestStorage.Reassign(pullRequest)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	updatedPRs, err := m.Storage.PullRequestStorage.Select(&pullRequestID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}
	if len(updatedPRs) == 0 {
		return domain.PullRequest{}, "", domain.ErrNotFound
	}
	updatedPullRequest := updatedPRs[0]

	return updatedPullRequest, newReviewer, nil
}

func (m *PullRequestManager) getReviewerTeamMembers(reviewerID string) ([]domain.TeamMember, error) {
	users, err := m.Storage.UserStorage.Select(&reviewerID)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrReviewerNotFound
	}
	reviewer := users[0]

	reviewerTeamName := reviewer.TeamName
	teams, err := m.Storage.TeamStorage.Select(&reviewerTeamName)
	if err != nil {
		return nil, err
	}
	if len(teams) == 0 {
		return nil, domain.ErrTeamNotFound
	}
	reviewerTeam := teams[0]

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

func (m *PullRequestManager) GetPullRequest(pullRequestID *string) (domain.PullRequest, error) {
	prs, err := m.Storage.PullRequestStorage.Select(pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}
	if len(prs) == 0 {
		return domain.PullRequest{}, domain.ErrNotFound
	}
	return prs[0], nil
}

func (m *PullRequestManager) GetPullRequests(pullRequestID *string) ([]domain.PullRequest, error) {
	return m.Storage.PullRequestStorage.Select(pullRequestID)
}
