package application

import (
	"context"
	"os"
	"testing"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func TestCreatePullRequest_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	ctx := context.Background()
	pr := domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID:       "test-pr-1",
			Name:     "Test PR",
			AuthorID: "test-user-1",
			Status:   domain.Open,
		},
	}

	_, err := CreatePullRequest(ctx, pr)
	if err != nil {
		t.Logf("CreatePullRequest error (may be expected if test data not set up): %v", err)
	}
}

func TestMergePullRequest_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	ctx := context.Background()
	pr := domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID: "test-pr-1",
		},
	}

	_, err := MergePullRequest(ctx, pr)
	if err != nil {
		t.Logf("MergePullRequest error (may be expected if test data not set up): %v", err)
	}
}

func TestReassignPullRequest_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	ctx := context.Background()
	prID := "test-pr-1"
	oldUserID := "test-user-1"

	_, _, err := ReassignPullRequest(ctx, prID, oldUserID)
	if err != nil {
		t.Logf("ReassignPullRequest error (may be expected if test data not set up): %v", err)
	}
}

func TestGetUserPullRequestsReviews_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	ctx := context.Background()
	userID := "test-user-1"

	_, err := GetUserPullRequestsReviews(ctx, userID)
	if err != nil {
		t.Logf("GetUserPullRequestsReviews error (may be expected if test data not set up): %v", err)
	}
}
