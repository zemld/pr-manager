package handlers

import (
	"strings"
	"time"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []domain.TeamMember `json:"members"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type PullRequestResponse struct {
	domain.PullRequestShort
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShortResponse struct {
	domain.PullRequestShort
}

type TeamWrapperResponse struct {
	Team domain.Team `json:"team"`
}

type UserWrapperResponse struct {
	User domain.User `json:"user"`
}

type PullRequestWrapperResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type ReassignResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type UserPullRequestsResponse struct {
	UserID       string                     `json:"user_id"`
	PullRequests []PullRequestShortResponse `json:"pull_requests"`
}

// Conversion functions
func domainPRToResponse(pr domain.PullRequest) PullRequestResponse {
	reviewersStr := strings.Trim(pr.AssignedReviewers, "[]")
	var reviewers []string
	if reviewersStr != "" {
		parts := strings.Split(reviewersStr, ",")
		for _, part := range parts {
			reviewer := strings.TrimSpace(part)
			if reviewer != "" {
				reviewers = append(reviewers, reviewer)
			}
		}
	}

	status := domain.PullRequestStatus(strings.ToUpper(string(pr.Status)))

	var createdAt *time.Time
	if pr.CreatedAt != nil {
		createdAt = pr.CreatedAt
	}

	var mergedAt *time.Time
	if pr.MergedAt != nil {
		mergedAt = pr.MergedAt
	}

	return PullRequestResponse{
		PullRequestShort: domain.PullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   status,
		},
		AssignedReviewers: reviewers,
		CreatedAt:         createdAt,
		MergedAt:          mergedAt,
	}
}

func domainPRToShortResponse(pr domain.PullRequest) PullRequestShortResponse {
	status := domain.PullRequestStatus(strings.ToUpper(string(pr.Status)))
	return PullRequestShortResponse{
		PullRequestShort: domain.PullRequestShort{
			ID:       pr.ID,
			Name:     pr.Name,
			AuthorID: pr.AuthorID,
			Status:   status,
		},
	}
}

func requestToDomainTeam(req CreateTeamRequest) domain.Team {
	return domain.Team{
		TeamName: req.TeamName,
		Members:  req.Members,
	}
}

func requestToDomainPR(req CreatePullRequestRequest) domain.PullRequest {
	return domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID:       req.PullRequestID,
			Name:     req.PullRequestName,
			AuthorID: req.AuthorID,
			Status:   domain.Open,
		},
	}
}

func requestToDomainPRForMerge(req MergePullRequestRequest) domain.PullRequest {
	return domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID: req.PullRequestID,
		},
	}
}
