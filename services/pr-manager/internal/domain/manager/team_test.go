package manager

import (
	"testing"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func TestTeamManager_AddTeam(t *testing.T) {
	tests := []struct {
		name     string
		team     domain.Team
		validate func(*testing.T, domain.Team, *mockTeamStorage)
	}{
		{
			name: "successfully add team with multiple members",
			team: createTestTeam("team1", []domain.TeamMember{
				{UserID: "user1", Username: "member1", IsActive: true},
				{UserID: "user2", Username: "member2", IsActive: true},
				{UserID: "user3", Username: "member3", IsActive: false},
			}),
			validate: func(t *testing.T, team domain.Team, storage *mockTeamStorage) {
				if team.TeamName != "team1" {
					t.Errorf("expected team name team1, got %s", team.TeamName)
				}
				if len(team.Members) != 3 {
					t.Errorf("expected 3 members, got %d", len(team.Members))
				}
				// Verify team was stored
				storedTeam, err := storage.Select("team1")
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if storedTeam.TeamName != "team1" {
					t.Errorf("expected stored team name team1, got %s", storedTeam.TeamName)
				}
			},
		},
		{
			name: "successfully add team with single member",
			team: createTestTeam("team2", []domain.TeamMember{
				{UserID: "user1", Username: "member1", IsActive: true},
			}),
			validate: func(t *testing.T, team domain.Team, storage *mockTeamStorage) {
				if len(team.Members) != 1 {
					t.Errorf("expected 1 member, got %d", len(team.Members))
				}
			},
		},
		{
			name: "successfully add empty team",
			team: createTestTeam("team3", []domain.TeamMember{}),
			validate: func(t *testing.T, team domain.Team, storage *mockTeamStorage) {
				if len(team.Members) != 0 {
					t.Errorf("expected 0 members, got %d", len(team.Members))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockTeamStorage()

			manager := NewTeamManager(storage)
			result, err := manager.AddTeam(tt.team)

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

func TestTeamManager_GetTeam(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockTeamStorage)
		teamName string
		wantErr  bool
		validate func(*testing.T, domain.Team)
	}{
		{
			name: "successfully get existing team",
			setup: func(storage *mockTeamStorage) {
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
					{UserID: "user2", Username: "member2", IsActive: true},
				})
				storage.Insert(team)
			},
			teamName: "team1",
			validate: func(t *testing.T, team domain.Team) {
				if team.TeamName != "team1" {
					t.Errorf("expected team name team1, got %s", team.TeamName)
				}
				if len(team.Members) != 2 {
					t.Errorf("expected 2 members, got %d", len(team.Members))
				}
				if team.Members[0].UserID != "user1" {
					t.Errorf("expected first member user1, got %s", team.Members[0].UserID)
				}
				if team.Members[1].UserID != "user2" {
					t.Errorf("expected second member user2, got %s", team.Members[1].UserID)
				}
			},
		},
		{
			name: "fail to get non-existent team",
			setup: func(storage *mockTeamStorage) {
				// No teams in storage
			},
			teamName: "nonexistent",
			wantErr:  true,
		},
		{
			name: "get team with inactive members",
			setup: func(storage *mockTeamStorage) {
				team := createTestTeam("team2", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
					{UserID: "user2", Username: "member2", IsActive: false},
					{UserID: "user3", Username: "member3", IsActive: true},
				})
				storage.Insert(team)
			},
			teamName: "team2",
			validate: func(t *testing.T, team domain.Team) {
				if len(team.Members) != 3 {
					t.Errorf("expected 3 members, got %d", len(team.Members))
				}
				// Verify inactive member is included
				if team.Members[1].UserID != "user2" || team.Members[1].IsActive {
					t.Error("expected user2 to be inactive")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := newMockTeamStorage()
			tt.setup(storage)

			manager := NewTeamManager(storage)
			result, err := manager.GetTeam(tt.teamName)

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
