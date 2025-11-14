package manager

import (
	"testing"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func TestUserManager_SelectUser(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockUserStorage)
		userID   string
		wantErr  bool
		validate func(*testing.T, domain.User)
	}{
		{
			name: "successfully select existing user",
			setup: func(storage *mockUserStorage) {
				user := createTestUser("user1", "testuser", "team1", true)
				storage.Insert(user)
			},
			userID: "user1",
			validate: func(t *testing.T, user domain.User) {
				if user.UserID != "user1" {
					t.Errorf("expected user ID user1, got %s", user.UserID)
				}
				if user.Username != "testuser" {
					t.Errorf("expected username testuser, got %s", user.Username)
				}
				if user.TeamName != "team1" {
					t.Errorf("expected team name team1, got %s", user.TeamName)
				}
				if !user.IsActive {
					t.Error("expected user to be active")
				}
			},
		},
		{
			name: "fail to select non-existent user",
			setup: func(storage *mockUserStorage) {
				// No users in storage
			},
			userID:  "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockUserStorage()
			tt.setup(storage)

			manager := NewUserManager(storage)
			result, err := manager.SelectUser(tt.userID)

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

func TestUserManager_UpdateUserStatus(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockUserStorage)
		user     domain.User
		wantErr  bool
		validate func(*testing.T, domain.User, *mockUserStorage)
	}{
		{
			name: "successfully update user status to inactive",
			setup: func(storage *mockUserStorage) {
				user := createTestUser("user1", "testuser", "team1", true)
				storage.Insert(user)
			},
			user: createTestUser("user1", "testuser", "team1", false),
			validate: func(t *testing.T, user domain.User, storage *mockUserStorage) {
				// Check that user was updated in storage
				updatedUser, err := storage.Select("user1")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if updatedUser.IsActive {
					t.Error("expected user to be inactive after update")
				}
			},
		},
		{
			name: "successfully update user status to active",
			setup: func(storage *mockUserStorage) {
				user := createTestUser("user1", "testuser", "team1", false)
				storage.Insert(user)
			},
			user: createTestUser("user1", "testuser", "team1", true),
			validate: func(t *testing.T, user domain.User, storage *mockUserStorage) {
				// Check that user was updated in storage
				updatedUser, err := storage.Select("user1")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !updatedUser.IsActive {
					t.Error("expected user to be active after update")
				}
			},
		},
		{
			name: "fail to update non-existent user",
			setup: func(storage *mockUserStorage) {
				// No users in storage
			},
			user:    createTestUser("nonexistent", "testuser", "team1", true),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockUserStorage()
			tt.setup(storage)

			manager := NewUserManager(storage)
			result, err := manager.UpdateUserStatus(tt.user)

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
				tt.validate(t, result, storage)
			}
		})
	}
}

func TestUserManager_UpdateUserStatus_SelectsUserFirst(t *testing.T) {
	// This test verifies that UpdateUserStatus calls SelectUser first
	// and only updates IsActive, preserving other fields
	storage := newMockUserStorage()
	user := createTestUser("user1", "testuser", "team1", true)
	storage.Insert(user)

	manager := NewUserManager(storage)

	// Try to update with different username (should preserve existing username)
	updateUser := createTestUser("user1", "differentname", "team1", false)
	result, err := manager.UpdateUserStatus(updateUser)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// The result should have the original username, not the update username
	// because SelectUser is called first and only IsActive is updated
	if result.Username != "testuser" {
		t.Errorf("expected username testuser (from SelectUser), got %s", result.Username)
	}

	// But the status should be updated
	if result.IsActive {
		t.Error("expected user to be inactive after update")
	}

	// Verify in storage
	updatedUser, _ := storage.Select("user1")
	if updatedUser.IsActive {
		t.Error("expected user to be inactive after update in storage")
	}
	if updatedUser.Username != "testuser" {
		t.Error("expected username to be preserved")
	}
}
