package manager

import (
	"errors"
	"strings"
	"time"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

var errNotFound = errors.New("not found")

// mockUserStorage is a mock implementation of storager.UserStorager
type mockUserStorage struct {
	users map[string]domain.User
}

func newMockUserStorage() *mockUserStorage {
	return &mockUserStorage{
		users: make(map[string]domain.User),
	}
}

func (m *mockUserStorage) Select(userID *string) ([]domain.User, error) {
	if userID == nil {
		// Return all users
		users := make([]domain.User, 0, len(m.users))
		for _, user := range m.users {
			users = append(users, user)
		}
		return users, nil
	}
	// Return specific user
	user, ok := m.users[*userID]
	if !ok {
		return []domain.User{}, nil
	}
	return []domain.User{user}, nil
}

func (m *mockUserStorage) Update(user domain.User) error {
	if _, ok := m.users[user.UserID]; !ok {
		return errNotFound
	}
	m.users[user.UserID] = user
	return nil
}

func (m *mockUserStorage) Insert(user domain.User) error {
	m.users[user.UserID] = user
	return nil
}

// mockTeamStorage is a mock implementation of storager.TeamStorager
type mockTeamStorage struct {
	teams map[string]domain.Team
}

func newMockTeamStorage() *mockTeamStorage {
	return &mockTeamStorage{
		teams: make(map[string]domain.Team),
	}
}

func (m *mockTeamStorage) Select(teamName *string) ([]domain.Team, error) {
	if teamName == nil {
		// Return all teams
		teams := make([]domain.Team, 0, len(m.teams))
		for _, team := range m.teams {
			teams = append(teams, team)
		}
		return teams, nil
	}
	// Return specific team
	team, ok := m.teams[*teamName]
	if !ok {
		return []domain.Team{}, nil
	}
	return []domain.Team{team}, nil
}

func (m *mockTeamStorage) Insert(team domain.Team) error {
	m.teams[team.TeamName] = team
	return nil
}

func (m *mockTeamStorage) Delete(teamName string) error {
	if _, ok := m.teams[teamName]; !ok {
		return errNotFound
	}
	delete(m.teams, teamName)
	return nil
}

// mockPullRequestStorage is a mock implementation of storager.PullRequestStorager
type mockPullRequestStorage struct {
	prs map[string]domain.PullRequest
}

func newMockPullRequestStorage() *mockPullRequestStorage {
	return &mockPullRequestStorage{
		prs: make(map[string]domain.PullRequest),
	}
}

func (m *mockPullRequestStorage) Select(pullRequestID *string) ([]domain.PullRequest, error) {
	if pullRequestID == nil {
		// Return all PRs
		prs := make([]domain.PullRequest, 0, len(m.prs))
		for _, pr := range m.prs {
			prs = append(prs, pr)
		}
		return prs, nil
	}
	// Return specific PR
	pr, ok := m.prs[*pullRequestID]
	if !ok {
		return []domain.PullRequest{}, nil
	}
	return []domain.PullRequest{pr}, nil
}

func (m *mockPullRequestStorage) Create(pullRequest domain.PullRequest) error {
	// Simulate ON CONFLICT DO NOTHING - if PR already exists, return error
	if _, exists := m.prs[pullRequest.ID]; exists {
		return errors.New("PR id already exists")
	}
	m.prs[pullRequest.ID] = pullRequest
	return nil
}

func (m *mockPullRequestStorage) Merge(pullRequest domain.PullRequest) error {
	pr, ok := m.prs[pullRequest.ID]
	if !ok {
		return errNotFound
	}
	// Simulate idempotent merge - only update if not already merged
	if pr.Status != domain.Merged {
		pr.Status = domain.Merged
		now := time.Now()
		pr.MergedAt = &now
		m.prs[pullRequest.ID] = pr
	}
	// If already merged, do nothing (idempotent)
	return nil
}

func (m *mockPullRequestStorage) Reassign(pullRequest domain.PullRequest) error {
	if _, ok := m.prs[pullRequest.ID]; !ok {
		return errNotFound
	}
	m.prs[pullRequest.ID] = pullRequest
	return nil
}

func (m *mockPullRequestStorage) SelectUserPullRequestsReviews(userID string) ([]domain.PullRequest, error) {
	var result []domain.PullRequest
	for _, pr := range m.prs {
		// Check if userID is in assigned_reviewers
		// Format is "[id1, id2]" or "[id1]"
		reviewers := strings.Trim(pr.AssignedReviewers, "[]")
		if reviewers != "" {
			parts := strings.Split(reviewers, ",")
			for _, part := range parts {
				if strings.TrimSpace(part) == userID {
					result = append(result, pr)
					break
				}
			}
		}
	}
	return result, nil
}

// createMockStorage creates a mock storage with all storages
func createMockStorage() *db.Storage {
	userStorage := newMockUserStorage()
	teamStorage := newMockTeamStorage()
	prStorage := newMockPullRequestStorage()

	return &db.Storage{
		UserStorage:        userStorage,
		TeamStorage:        teamStorage,
		PullRequestStorage: prStorage,
	}
}

// Helper function to create test users
func createTestUser(userID, username, teamName string, isActive bool) domain.User {
	return domain.User{
		UserID:   userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}
}

// Helper function to create test team
func createTestTeam(teamName string, members []domain.TeamMember) domain.Team {
	return domain.Team{
		TeamName: teamName,
		Members:  members,
	}
}

// Helper function to create test PR
func createTestPR(id, name, authorID string, status domain.PullRequestStatus, reviewers string) domain.PullRequest {
	return domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID:       id,
			Name:     name,
			AuthorID: authorID,
			Status:   status,
		},
		AssignedReviewers: reviewers,
		CreatedAt:         func() *time.Time { t := time.Now(); return &t }(),
		MergedAt:          nil,
	}
}
