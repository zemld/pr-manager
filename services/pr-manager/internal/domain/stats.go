package domain

type Stats struct {
	UserStats           UserStats                      `json:"user_stats"`
	IndividualUserStats map[string]IndividualUserStats `json:"individual_user_stats"`

	TeamStats           TeamStats                      `json:"team_stats"`
	IndividualTeamStats map[string]IndividualTeamStats `json:"individual_team_stats"`

	PullRequestStats PullRequestStats `json:"pull_request_stats"`
}

type UserStats struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Inactive int64 `json:"inactive"`
}

type IndividualUserStats struct {
	Username              string  `json:"username"`
	PRsCreated            int64   `json:"prs_created"`
	PRsReviewed           int64   `json:"prs_reviewed"`
	PRsMerged             int64   `json:"prs_merged"`
	PRsOpen               int64   `json:"prs_open"`
	PRsWaitingForReview   int64   `json:"prs_waiting_for_review"`
	AverageMergeTimeHours float64 `json:"average_merge_time_hours"`
}

type TeamStats struct {
	Total                         int64   `json:"total"`
	AverageMembersPerTeam         float64 `json:"average_members_per_team"`
	MostMembersInTeam             int64   `json:"most_members_in_team"`
	LeastMembersInTeam            int64   `json:"least_members_in_team"`
	AverageActiveMembersPerTeam   float64 `json:"average_active_members_per_team"`
	MostActiveMembersInTeam       int64   `json:"most_active_members_in_team"`
	LeastActiveMembersInTeam      int64   `json:"least_active_members_in_team"`
	AverageInactiveMembersPerTeam float64 `json:"average_inactive_members_per_team"`
	MostInactiveMembersInTeam     int64   `json:"most_inactive_members_in_team"`
	LeastInactiveMembersInTeam    int64   `json:"least_inactive_members_in_team"`
}

type IndividualTeamStats struct {
	TotalMembers          int64   `json:"total_members"`
	ActiveMembers         int64   `json:"active_members"`
	InactiveMembers       int64   `json:"inactive_members"`
	PRsCreated            int64   `json:"prs_created"`
	PRsReviewed           int64   `json:"prs_reviewed"`
	PRsMerged             int64   `json:"prs_merged"`
	PRsOpen               int64   `json:"prs_open"`
	PRsWaitingForReview   int64   `json:"prs_waiting_for_review"`
	AverageMergeTimeHours float64 `json:"average_merge_time_hours"`
}

type PullRequestStats struct {
	Total                 int64   `json:"total"`
	AveragePRsPerUser     float64 `json:"average_prs_per_user"`
	MostPRsPerUser        int64   `json:"most_prs_per_user"`
	LeastPRsPerUser       int64   `json:"least_prs_per_user"`
	AverageMergeTimeHours float64 `json:"average_merge_time_hours"`
	AveragePRsPerReviewer float64 `json:"average_prs_per_reviewer"`
	MostPRsPerReviewer    int64   `json:"most_prs_per_reviewer"`
	LeastPRsPerReviewer   int64   `json:"least_prs_per_reviewer"`
}
