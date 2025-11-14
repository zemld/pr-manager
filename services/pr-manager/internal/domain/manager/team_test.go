package manager

import (
	"strings"
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
				teamName := "team1"
				storedTeams, err := storage.Select(&teamName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(storedTeams) == 0 {
					t.Error("expected team to be found")
					return
				}
				storedTeam := storedTeams[0]
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

			manager := NewTeamManager(storage, nil)
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

			manager := NewTeamManager(storage, nil)
			teamNamePtr := &tt.teamName
			result, err := manager.GetTeam(teamNamePtr)

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

func TestTeamManager_DeleteTeam(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockTeamStorage, *mockPullRequestStorage)
		teamName string
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *mockTeamStorage, *mockPullRequestStorage)
	}{
		{
			name: "successfully delete team without PRs",
			setup: func(teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
					{UserID: "user2", Username: "member2", IsActive: true},
				})
				teamStorage.Insert(team)
			},
			teamName: "team1",
			validate: func(t *testing.T, teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				teamName := "team1"
				teams, err := teamStorage.Select(&teamName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(teams) != 0 {
					t.Error("expected team to be deleted")
				}
			},
		},
		{
			name: "successfully delete team and remove reviewers from PRs",
			setup: func(teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
					{UserID: "user2", Username: "member2", IsActive: true},
				})
				teamStorage.Insert(team)

				// Create PRs with reviewers from team1
				pr1 := createTestPR("pr1", "PR 1", "author1", domain.Open, "[user1, user2]")
				pr2 := createTestPR("pr2", "PR 2", "author2", domain.Open, "[user1, other_reviewer]")
				pr3 := createTestPR("pr3", "PR 3", "author3", domain.Open, "[other_reviewer, another_reviewer]")
				prStorage.Create(pr1)
				prStorage.Create(pr2)
				prStorage.Create(pr3)
			},
			teamName: "team1",
			validate: func(t *testing.T, teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				// Verify team is deleted
				teamName := "team1"
				teams, err := teamStorage.Select(&teamName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(teams) != 0 {
					t.Error("expected team to be deleted")
				}

				// Verify reviewers are removed from PRs
				prID1 := "pr1"
				prs1, err := prStorage.Select(&prID1)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(prs1) == 0 {
					t.Error("expected pr1 to exist")
					return
				}
				if prs1[0].AssignedReviewers != "[]" {
					t.Errorf("expected pr1 to have empty reviewers, got %s", prs1[0].AssignedReviewers)
				}

				prID2 := "pr2"
				prs2, err := prStorage.Select(&prID2)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(prs2) == 0 {
					t.Error("expected pr2 to exist")
					return
				}
				// pr2 should have only other_reviewer, user1 should be removed
				if strings.Contains(prs2[0].AssignedReviewers, "user1") {
					t.Error("expected user1 to be removed from pr2 reviewers")
				}
				if !strings.Contains(prs2[0].AssignedReviewers, "other_reviewer") {
					t.Error("expected other_reviewer to remain in pr2 reviewers")
				}

				// pr3 should remain unchanged
				prID3 := "pr3"
				prs3, err := prStorage.Select(&prID3)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(prs3) == 0 {
					t.Error("expected pr3 to exist")
					return
				}
				if prs3[0].AssignedReviewers != "[other_reviewer, another_reviewer]" {
					t.Errorf("expected pr3 reviewers to remain unchanged, got %s", prs3[0].AssignedReviewers)
				}
			},
		},
		{
			name: "successfully delete team with reviewers in merged PRs",
			setup: func(teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
				})
				teamStorage.Insert(team)

				// Create merged PR with reviewer from team1
				pr1 := createTestPR("pr1", "PR 1", "author1", domain.Merged, "[user1, user2]")
				prStorage.Create(pr1)
				// Merge it
				prStorage.Merge(pr1)
			},
			teamName: "team1",
			validate: func(t *testing.T, teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				// Verify team is deleted
				teamName := "team1"
				teams, err := teamStorage.Select(&teamName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(teams) != 0 {
					t.Error("expected team to be deleted")
				}

				// Verify reviewers are removed from merged PR too
				prID1 := "pr1"
				prs1, err := prStorage.Select(&prID1)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(prs1) == 0 {
					t.Error("expected pr1 to exist")
					return
				}
				// user1 should be removed, user2 should remain
				if strings.Contains(prs1[0].AssignedReviewers, "user1") {
					t.Error("expected user1 to be removed from pr1 reviewers")
				}
				if !strings.Contains(prs1[0].AssignedReviewers, "user2") {
					t.Error("expected user2 to remain in pr1 reviewers")
				}
			},
		},
		{
			name: "fail to delete non-existent team",
			setup: func(teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				// No teams in storage
			},
			teamName: "nonexistent",
			wantErr:  true,
			errMsg:   "team not found",
		},
		{
			name: "successfully delete team with multiple reviewers in same PR",
			setup: func(teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				team := createTestTeam("team1", []domain.TeamMember{
					{UserID: "user1", Username: "member1", IsActive: true},
					{UserID: "user2", Username: "member2", IsActive: true},
				})
				teamStorage.Insert(team)

				// Create PR with both team members as reviewers
				pr1 := createTestPR("pr1", "PR 1", "author1", domain.Open, "[user1, user2, other_reviewer]")
				prStorage.Create(pr1)
			},
			teamName: "team1",
			validate: func(t *testing.T, teamStorage *mockTeamStorage, prStorage *mockPullRequestStorage) {
				// Verify team is deleted
				teamName := "team1"
				teams, err := teamStorage.Select(&teamName)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(teams) != 0 {
					t.Error("expected team to be deleted")
				}

				// Verify both team members are removed, other_reviewer remains
				prID1 := "pr1"
				prs1, err := prStorage.Select(&prID1)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(prs1) == 0 {
					t.Error("expected pr1 to exist")
					return
				}
				if strings.Contains(prs1[0].AssignedReviewers, "user1") {
					t.Error("expected user1 to be removed from pr1 reviewers")
				}
				if strings.Contains(prs1[0].AssignedReviewers, "user2") {
					t.Error("expected user2 to be removed from pr1 reviewers")
				}
				if !strings.Contains(prs1[0].AssignedReviewers, "other_reviewer") {
					t.Error("expected other_reviewer to remain in pr1 reviewers")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teamStorage := newMockTeamStorage()
			prStorage := newMockPullRequestStorage()
			tt.setup(teamStorage, prStorage)

			manager := NewTeamManager(teamStorage, prStorage)
			err := manager.DeleteTeam(tt.teamName)

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
				tt.validate(t, teamStorage, prStorage)
			}
		})
	}
}
