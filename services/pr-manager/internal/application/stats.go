package application

import (
	"context"
	"strings"
	"sync"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type data struct {
	mu           sync.Mutex
	users        []domain.User
	teams        []domain.Team
	pullRequests []domain.PullRequest
}

func GetStats(ctx context.Context) domain.Stats {
	d := data{}
	wg := sync.WaitGroup{}
	wg.Add(3)

	go fillUsersData(ctx, &d, &wg)
	go fillTeamsData(ctx, &d, &wg)
	go fillPullRequestsData(ctx, &d, &wg)

	wg.Wait()

	stats := domain.Stats{}
	calculateUserStats(&d, &stats)
	calculateIndividualUserStats(&d, &stats)
	calculateTeamStats(&d, &stats)
	calculateIndividualTeamStats(&d, &stats)
	calculatePullRequestStats(&d, &stats)

	return stats
}

func fillUsersData(ctx context.Context, d *data, wg *sync.WaitGroup) {
	defer wg.Done()

	users, err := GetUsers(ctx)
	d.mu.Lock()
	defer d.mu.Unlock()
	if err == nil {
		d.users = users
	}
}

func fillTeamsData(ctx context.Context, d *data, wg *sync.WaitGroup) {
	defer wg.Done()

	teams, err := GetTeams(ctx)
	d.mu.Lock()
	defer d.mu.Unlock()
	if err == nil {
		d.teams = teams
	}
}

func fillPullRequestsData(ctx context.Context, d *data, wg *sync.WaitGroup) {
	defer wg.Done()

	pullRequests, err := GetPullRequests(ctx)
	d.mu.Lock()
	defer d.mu.Unlock()
	if err == nil {
		d.pullRequests = pullRequests
	}
}

func calculateUserStats(d *data, stats *domain.Stats) {
	stats.UserStats.Total = int64(len(d.users))
	activeUsers := int64(0)
	inactiveUsers := int64(0)
	for _, user := range d.users {
		if user.IsActive {
			activeUsers++
		} else {
			inactiveUsers++
		}
	}
	stats.UserStats.Active = activeUsers
	stats.UserStats.Inactive = inactiveUsers
}

func calculateIndividualUserStats(d *data, stats *domain.Stats) {

	individualUserStats := make(map[string]domain.IndividualUserStats, len(d.users))
	for _, user := range d.users {
		prsCreated := getFilteredPRsForUser(user, d.pullRequests, func(pr domain.PullRequest, user domain.User) bool {
			return pr.AuthorID == user.UserID
		})
		prsReviewed := getFilteredPRsForUser(user, d.pullRequests, func(pr domain.PullRequest, user domain.User) bool {
			return strings.Contains(pr.AssignedReviewers, user.UserID) && pr.Status == domain.Merged
		})
		prsMerged := getFilteredPRsForUser(user, d.pullRequests, func(pr domain.PullRequest, user domain.User) bool {
			return pr.Status == domain.Merged && pr.AuthorID == user.UserID
		})
		prsOpen := getFilteredPRsForUser(user, d.pullRequests, func(pr domain.PullRequest, user domain.User) bool {
			return pr.Status == domain.Open && pr.AuthorID == user.UserID
		})
		prsWaitingForReview := getFilteredPRsForUser(user, d.pullRequests, func(pr domain.PullRequest, user domain.User) bool {
			return strings.Contains(pr.AssignedReviewers, user.UserID) && pr.Status == domain.Open
		})
		averageMergeTimeHours := 0.0
		if len(prsMerged) > 0 {
			for _, pr := range prsMerged {
				averageMergeTimeHours += pr.MergedAt.Sub(*pr.CreatedAt).Hours()
			}
			averageMergeTimeHours /= float64(len(prsMerged))
		}
		individualUserStats[user.UserID] = domain.IndividualUserStats{
			Username:              user.Username,
			PRsCreated:            int64(len(prsCreated)),
			PRsReviewed:           int64(len(prsReviewed)),
			PRsMerged:             int64(len(prsMerged)),
			PRsOpen:               int64(len(prsOpen)),
			PRsWaitingForReview:   int64(len(prsWaitingForReview)),
			AverageMergeTimeHours: averageMergeTimeHours,
		}
	}
	stats.IndividualUserStats = individualUserStats
}
func getFilteredPRsForUser(user domain.User, prs []domain.PullRequest, filter func(domain.PullRequest, domain.User) bool) []domain.PullRequest {
	filteredPRs := make([]domain.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if filter(pr, user) {
			filteredPRs = append(filteredPRs, pr)
		}
	}
	return filteredPRs
}

func calculateTeamStats(d *data, stats *domain.Stats) {
	stats.TeamStats.Total = int64(len(d.teams))
	if stats.TeamStats.Total > 0 {
		stats.TeamStats.AverageMembersPerTeam = float64(stats.UserStats.Total) / float64(stats.TeamStats.Total)
		stats.TeamStats.AverageActiveMembersPerTeam = float64(stats.UserStats.Active) / float64(stats.TeamStats.Total)
		stats.TeamStats.AverageInactiveMembersPerTeam = float64(stats.UserStats.Inactive) / float64(stats.TeamStats.Total)
	}
	if len(d.teams) > 0 {
		stats.TeamStats.LeastMembersInTeam = int64(len(d.teams[0].Members))
		stats.TeamStats.LeastActiveMembersInTeam = int64(len(d.teams[0].Members))
		stats.TeamStats.LeastInactiveMembersInTeam = int64(len(d.teams[0].Members))
	}
	for _, team := range d.teams {
		if int64(len(team.Members)) > stats.TeamStats.MostMembersInTeam {
			stats.TeamStats.MostMembersInTeam = int64(len(team.Members))
		}
		if int64(len(team.Members)) < stats.TeamStats.LeastMembersInTeam {
			stats.TeamStats.LeastMembersInTeam = int64(len(team.Members))
		}
		activeMembers := 0
		inactiveMembers := 0
		for _, member := range team.Members {
			if member.IsActive {
				activeMembers++
			} else {
				inactiveMembers++
			}
		}
		if int64(activeMembers) > stats.TeamStats.MostActiveMembersInTeam {
			stats.TeamStats.MostActiveMembersInTeam = int64(activeMembers)
		}
		if int64(inactiveMembers) > stats.TeamStats.MostInactiveMembersInTeam {
			stats.TeamStats.MostInactiveMembersInTeam = int64(inactiveMembers)
		}
		if int64(inactiveMembers) < stats.TeamStats.LeastInactiveMembersInTeam {
			stats.TeamStats.LeastInactiveMembersInTeam = int64(inactiveMembers)
		}
	}
}

func calculateIndividualTeamStats(d *data, stats *domain.Stats) {
	individualTeamStats := make(map[string]domain.IndividualTeamStats, len(d.teams))
	for _, team := range d.teams {
		activeMembers := 0
		inactiveMembers := 0
		prsCreated := int64(0)
		prsReviewed := int64(0)
		prsMerged := int64(0)
		prsOpen := int64(0)
		prsWaitingForReview := int64(0)
		averageMergeTimeHours := 0.0
		for _, member := range team.Members {
			if member.IsActive {
				activeMembers++
			} else {
				inactiveMembers++
			}
			prsCreated += stats.IndividualUserStats[member.UserID].PRsCreated
			prsReviewed += stats.IndividualUserStats[member.UserID].PRsReviewed
			prsMerged += stats.IndividualUserStats[member.UserID].PRsMerged
			prsOpen += stats.IndividualUserStats[member.UserID].PRsOpen
			prsWaitingForReview += stats.IndividualUserStats[member.UserID].PRsWaitingForReview
			averageMergeTimeHours += stats.IndividualUserStats[member.UserID].AverageMergeTimeHours
		}
		if len(team.Members) > 0 {
			averageMergeTimeHours /= float64(len(team.Members))
		}
		individualTeamStats[team.TeamName] = domain.IndividualTeamStats{
			TotalMembers:          int64(len(team.Members)),
			ActiveMembers:         int64(activeMembers),
			InactiveMembers:       int64(inactiveMembers),
			PRsCreated:            prsCreated,
			PRsReviewed:           prsReviewed,
			PRsMerged:             prsMerged,
			PRsOpen:               prsOpen,
			PRsWaitingForReview:   prsWaitingForReview,
			AverageMergeTimeHours: averageMergeTimeHours,
		}
	}
	stats.IndividualTeamStats = individualTeamStats
}

func calculatePullRequestStats(d *data, stats *domain.Stats) {
	stats.PullRequestStats.Total = int64(len(d.pullRequests))
	reviewersCount := int64(0)
	prCreatorsCount := int64(0)
	averageMergeTimeHours := 0.0
	if stats.UserStats.Total > 0 {
		stats.PullRequestStats.LeastPRsPerUser = stats.IndividualUserStats[d.users[0].UserID].PRsCreated
		stats.PullRequestStats.LeastPRsPerReviewer = stats.IndividualUserStats[d.users[0].UserID].PRsReviewed
	}
	if stats.UserStats.Total > 0 {
		stats.PullRequestStats.AveragePRsPerUser = float64(stats.PullRequestStats.Total) / float64(stats.UserStats.Total)
	}
	for _, userStats := range stats.IndividualUserStats {
		if userStats.PRsCreated > stats.PullRequestStats.MostPRsPerUser {
			stats.PullRequestStats.MostPRsPerUser = userStats.PRsCreated
		}
		if userStats.PRsCreated < stats.PullRequestStats.LeastPRsPerUser {
			stats.PullRequestStats.LeastPRsPerUser = userStats.PRsCreated
		}
		averageMergeTimeHours += userStats.AverageMergeTimeHours
		stats.PullRequestStats.AveragePRsPerReviewer += float64(userStats.PRsReviewed)
		if userStats.PRsReviewed > stats.PullRequestStats.MostPRsPerReviewer {
			stats.PullRequestStats.MostPRsPerReviewer = userStats.PRsReviewed
		}
		if userStats.PRsReviewed < stats.PullRequestStats.LeastPRsPerReviewer {
			stats.PullRequestStats.LeastPRsPerReviewer = userStats.PRsReviewed
		}
		if userStats.PRsReviewed > 0 {
			reviewersCount++
		}
		if userStats.PRsCreated > 0 {
			prCreatorsCount++
		}
	}
	if reviewersCount > 0 {
		stats.PullRequestStats.AveragePRsPerReviewer /= float64(reviewersCount)
	}
	if prCreatorsCount > 0 {
		stats.PullRequestStats.AverageMergeTimeHours = averageMergeTimeHours / float64(prCreatorsCount)
	}
}
