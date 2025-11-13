package domain

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type TeamMember struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
