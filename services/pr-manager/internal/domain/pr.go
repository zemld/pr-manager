package domain

import "time"

type PullRequest struct {
	PullRequestShort
	AssignedReviewers string    `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"created_at"`
	MergedAt          time.Time `json:"merged_at"`
}

type PullRequestShort struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	AuthorID string            `json:"author_id"`
	Status   PullRequestStatus `json:"status"`
}

type PullRequestStatus string

const (
	Open   PullRequestStatus = "open"
	Merged PullRequestStatus = "merged"
)
