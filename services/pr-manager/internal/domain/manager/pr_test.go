package manager

import (
	"strings"
	"testing"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
	"github.com/zemld/pr-manager/pr-manager/internal/domain/db"
)

func TestPullRequestManager_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*db.Storage)
		pr       domain.PullRequest
		wantErr  bool
		errMsg   string
		validate func(*testing.T, domain.PullRequest, *db.Storage)
	}{
		{
			name: "successfully create PR with 2 reviewers",
			setup: func(storage *db.Storage) {
				// Create team with 3 active members (author + 2 reviewers)
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
					{UserID: "user3", Username: "reviewer2", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				// Create users
				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", true))
			},
			pr: createTestPR("pr1", "Test PR", "user1", domain.Open, ""),
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				if pr.Status != domain.Open {
					t.Errorf("expected status Open, got %v", pr.Status)
				}
				if pr.AssignedReviewers == "" {
					t.Error("expected reviewers to be assigned")
				}
				// Check that reviewers are in the format [user2, user3] or [user3, user2]
				if !strings.HasPrefix(pr.AssignedReviewers, "[") || !strings.HasSuffix(pr.AssignedReviewers, "]") {
					t.Errorf("expected reviewers in format [id1, id2], got %s", pr.AssignedReviewers)
				}
				// Check that author is not in reviewers
				if strings.Contains(pr.AssignedReviewers, "user1") {
					t.Error("author should not be in reviewers")
				}
				// Check that we have 2 reviewers
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 2 {
					t.Errorf("expected 2 reviewers, got %d", len(reviewerList))
				}
			},
		},
		{
			name: "create PR with only 1 available reviewer",
			setup: func(storage *db.Storage) {
				// Create team with 2 active members (author + 1 reviewer)
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
			},
			pr: createTestPR("pr1", "Test PR", "user1", domain.Open, ""),
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				if pr.AssignedReviewers == "" {
					t.Error("expected at least one reviewer to be assigned")
				}
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 1 {
					t.Errorf("expected 1 reviewer, got %d", len(reviewerList))
				}
				if strings.Contains(pr.AssignedReviewers, "user1") {
					t.Error("author should not be in reviewers")
				}
			},
		},
		{
			name: "create PR with no available reviewers (only author in team)",
			setup: func(storage *db.Storage) {
				// Create team with only author
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
				})
				storage.TeamStorage.Insert(team)
				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
			},
			pr:      createTestPR("pr1", "Test PR", "user1", domain.Open, ""),
			wantErr: true,
			errMsg:  "no possible assigners",
		},
		{
			name: "create PR with inactive reviewers (should not assign inactive)",
			setup: func(storage *db.Storage) {
				// Create team with author and inactive reviewers
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: false},
					{UserID: "user3", Username: "reviewer2", IsActive: false},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", false))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", false))
			},
			pr:      createTestPR("pr1", "Test PR", "user1", domain.Open, ""),
			wantErr: true,
			errMsg:  "no possible assigners",
		},
		{
			name: "create PR with mix of active and inactive reviewers",
			setup: func(storage *db.Storage) {
				// Create team with author, 1 active and 1 inactive reviewer
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
					{UserID: "user3", Username: "reviewer2", IsActive: false},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", false))
			},
			pr: createTestPR("pr1", "Test PR", "user1", domain.Open, ""),
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				if pr.AssignedReviewers == "" {
					t.Error("expected at least one reviewer to be assigned")
				}
				// Should only assign active reviewer
				if !strings.Contains(pr.AssignedReviewers, "user2") {
					t.Error("expected user2 to be assigned")
				}
				if strings.Contains(pr.AssignedReviewers, "user3") {
					t.Error("inactive user3 should not be assigned")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := createMockStorage()
			tt.setup(storage)

			manager := NewPullRequestManager(storage)
			result, err := manager.CreatePullRequest(tt.pr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.pr.ID {
				t.Errorf("expected PR ID %s, got %s", tt.pr.ID, result.ID)
			}

			if tt.validate != nil {
				tt.validate(t, result, storage)
			}
		})
	}
}

func TestPullRequestManager_MergePullRequest(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*db.Storage)
		pr       domain.PullRequest
		wantErr  bool
		validate func(*testing.T, domain.PullRequest)
	}{
		{
			name: "successfully merge open PR",
			setup: func(storage *db.Storage) {
				pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2, user3]")
				storage.PullRequestStorage.Create(pr)
			},
			pr: createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2, user3]"),
			validate: func(t *testing.T, pr domain.PullRequest) {
				if pr.Status != domain.Merged {
					t.Errorf("expected status Merged, got %v", pr.Status)
				}
			},
		},
		{
			name: "merge already merged PR (idempotent)",
			setup: func(storage *db.Storage) {
				pr := createTestPR("pr1", "Test PR", "user1", domain.Merged, "[user2, user3]")
				storage.PullRequestStorage.Create(pr)
			},
			pr: createTestPR("pr1", "Test PR", "user1", domain.Merged, "[user2, user3]"),
			validate: func(t *testing.T, pr domain.PullRequest) {
				if pr.Status != domain.Merged {
					t.Errorf("expected status Merged, got %v", pr.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := createMockStorage()
			tt.setup(storage)

			manager := NewPullRequestManager(storage)
			result, err := manager.MergePullRequest(tt.pr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestPullRequestManager_ReassignPullRequest(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*db.Storage)
		prID          string
		oldReviewerID string
		wantErr       bool
		errMsg        string
		validate      func(*testing.T, domain.PullRequest, *db.Storage)
	}{
		{
			name: "successfully reassign reviewer",
			setup: func(storage *db.Storage) {
				// Create PR with 2 reviewers
				pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2, user3]")
				storage.PullRequestStorage.Create(pr)

				// Create team for old reviewer (user2)
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
					{UserID: "user3", Username: "reviewer2", IsActive: true},
					{UserID: "user4", Username: "reviewer3", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", true))
				storage.UserStorage.Insert(createTestUser("user4", "reviewer3", "team1", true))
			},
			prID:          "pr1",
			oldReviewerID: "user2",
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				if pr.Status != domain.Open {
					t.Errorf("expected status Open, got %v", pr.Status)
				}
				// Should not contain old reviewer
				if strings.Contains(pr.AssignedReviewers, "user2") {
					t.Error("old reviewer user2 should not be in reviewers")
				}
				// Should contain the other reviewer (user3)
				if !strings.Contains(pr.AssignedReviewers, "user3") {
					t.Error("expected user3 to remain in reviewers")
				}
				// Should have a new reviewer from user2's team
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 2 {
					t.Errorf("expected 2 reviewers, got %d", len(reviewerList))
				}
			},
		},
		{
			name: "reassign when only one reviewer exists",
			setup: func(storage *db.Storage) {
				// Create PR with 1 reviewer
				pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2]")
				storage.PullRequestStorage.Create(pr)

				// Create team for old reviewer (user2)
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
					{UserID: "user4", Username: "reviewer3", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user4", "reviewer3", "team1", true))
			},
			prID:          "pr1",
			oldReviewerID: "user2",
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				// Should not contain old reviewer
				if strings.Contains(pr.AssignedReviewers, "user2") {
					t.Error("old reviewer user2 should not be in reviewers")
				}
				// Should have a new reviewer
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 1 {
					t.Errorf("expected 1 reviewer, got %d", len(reviewerList))
				}
			},
		},
		{
			name: "fail to reassign merged PR",
			setup: func(storage *db.Storage) {
				// Create merged PR
				pr := createTestPR("pr1", "Test PR", "user1", domain.Merged, "[user2, user3]")
				pr.Status = domain.Merged
				storage.PullRequestStorage.Create(pr)
			},
			prID:          "pr1",
			oldReviewerID: "user2",
			wantErr:       true,
			errMsg:        "pull request is already merged",
		},
		{
			name: "reassign when no available replacement (only old reviewer in team)",
			setup: func(storage *db.Storage) {
				// Create PR with 2 reviewers
				pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2, user3]")
				storage.PullRequestStorage.Create(pr)

				// Create team with only old reviewer (no replacement available)
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", true))
			},
			prID:          "pr1",
			oldReviewerID: "user2",
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				// Should not contain old reviewer
				if strings.Contains(pr.AssignedReviewers, "user2") {
					t.Error("old reviewer user2 should not be in reviewers")
				}
				// Should only have the other reviewer (user3)
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 1 {
					t.Errorf("expected 1 reviewer, got %d", len(reviewerList))
				}
				if !strings.Contains(pr.AssignedReviewers, "user3") {
					t.Error("expected user3 to remain in reviewers")
				}
			},
		},
		{
			name: "reassign excludes author and other reviewer from new assignment",
			setup: func(storage *db.Storage) {
				// Create PR with 2 reviewers
				pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "[user2, user3]")
				storage.PullRequestStorage.Create(pr)

				// Create team with multiple members
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "author", IsActive: true},
					{UserID: "user2", Username: "reviewer1", IsActive: true},
					{UserID: "user3", Username: "reviewer2", IsActive: true},
					{UserID: "user4", Username: "reviewer3", IsActive: true},
					{UserID: "user5", Username: "reviewer4", IsActive: true},
				})
				storage.TeamStorage.Insert(team)

				storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
				storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
				storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", true))
				storage.UserStorage.Insert(createTestUser("user4", "reviewer3", "team1", true))
				storage.UserStorage.Insert(createTestUser("user5", "reviewer4", "team1", true))
			},
			prID:          "pr1",
			oldReviewerID: "user2",
			validate: func(t *testing.T, pr domain.PullRequest, storage *db.Storage) {
				// Should not contain old reviewer, author, or other reviewer
				if strings.Contains(pr.AssignedReviewers, "user2") {
					t.Error("old reviewer user2 should not be in reviewers")
				}
				if strings.Contains(pr.AssignedReviewers, "user1") {
					t.Error("author user1 should not be in reviewers")
				}
				// Should contain user3 (the other reviewer)
				if !strings.Contains(pr.AssignedReviewers, "user3") {
					t.Error("expected user3 to remain in reviewers")
				}
				// New reviewer should be user4 or user5 (not user1, user2, or user3)
				reviewers := strings.Trim(pr.AssignedReviewers, "[]")
				reviewerList := strings.Split(reviewers, ",")
				if len(reviewerList) != 2 {
					t.Errorf("expected 2 reviewers, got %d", len(reviewerList))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := createMockStorage()
			tt.setup(storage)

			manager := NewPullRequestManager(storage)
			result, err := manager.ReassignPullRequest(tt.prID, tt.oldReviewerID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result, storage)
			}
		})
	}
}

func TestPullRequestManager_HelperMethods(t *testing.T) {
	// Test helper methods indirectly through CreatePullRequest
	// since they are not exported

	t.Run("getActiveUserIDsFromTeam filters inactive users", func(t *testing.T) {
		storage := createMockStorage()

		// Create team with mix of active and inactive members
		team := createTestTeam("team1", []domain.TeamMember{
			{UserID: "user1", Username: "author", IsActive: true},
			{UserID: "user2", Username: "reviewer1", IsActive: true},
			{UserID: "user3", Username: "reviewer2", IsActive: false}, // inactive
			{UserID: "user4", Username: "reviewer3", IsActive: true},
		})
		storage.TeamStorage.Insert(team)

		storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
		storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
		storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", false))
		storage.UserStorage.Insert(createTestUser("user4", "reviewer3", "team1", true))

		manager := NewPullRequestManager(storage)
		pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "")
		result, err := manager.CreatePullRequest(pr)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Verify that inactive user3 is not assigned
		if strings.Contains(result.AssignedReviewers, "user3") {
			t.Error("inactive user3 should not be assigned")
		}
		// Verify that active users are assigned (but not author)
		if strings.Contains(result.AssignedReviewers, "user1") {
			t.Error("author user1 should not be assigned")
		}
	})

	t.Run("filterReviewers excludes specified users", func(t *testing.T) {
		storage := createMockStorage()

		// Create team with 4 members (author + 3 potential reviewers)
		team := createTestTeam("team1", []domain.TeamMember{
			{UserID: "user1", Username: "author", IsActive: true},
			{UserID: "user2", Username: "reviewer1", IsActive: true},
			{UserID: "user3", Username: "reviewer2", IsActive: true},
			{UserID: "user4", Username: "reviewer3", IsActive: true},
		})
		storage.TeamStorage.Insert(team)

		storage.UserStorage.Insert(createTestUser("user1", "author", "team1", true))
		storage.UserStorage.Insert(createTestUser("user2", "reviewer1", "team1", true))
		storage.UserStorage.Insert(createTestUser("user3", "reviewer2", "team1", true))
		storage.UserStorage.Insert(createTestUser("user4", "reviewer3", "team1", true))

		manager := NewPullRequestManager(storage)
		pr := createTestPR("pr1", "Test PR", "user1", domain.Open, "")
		result, err := manager.CreatePullRequest(pr)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		// Verify that author is filtered out
		if strings.Contains(result.AssignedReviewers, "user1") {
			t.Error("author user1 should be filtered out")
		}
		// Should assign 2 reviewers from user2, user3, user4
		reviewers := strings.Trim(result.AssignedReviewers, "[]")
		reviewerList := strings.Split(reviewers, ",")
		if len(reviewerList) != 2 {
			t.Errorf("expected 2 reviewers, got %d", len(reviewerList))
		}
	})
}

// Helper function to check if slice contains value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
