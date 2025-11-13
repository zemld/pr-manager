package domain

import "time"

type PullRequest struct {
	PullRequestShort
	AssignedReviewers string
	CreatedAt         time.Time
	MergedAt          time.Time
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   PullRequestStatus
}

type PullRequestStatus string

const (
	Open   PullRequestStatus = "open"
	Merged PullRequestStatus = "merged"
)
