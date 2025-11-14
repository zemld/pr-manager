package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func TestCreatePullRequestHandler_Integration(t *testing.T) {
	// Skip if database is not available
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	tests := []struct {
		name           string
		requestBody    CreatePullRequestRequest
		expectedStatus int
		expectedCode   ErrorCode
		skip           bool
	}{
		{
			name:           "invalid request body",
			requestBody:    CreatePullRequestRequest{},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorCodeNotFound,
		},
		{
			name: "missing required fields",
			requestBody: CreatePullRequestRequest{
				PullRequestID: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorCodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Test skipped")
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			CreatePullRequestHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err == nil {
					if errResp.Error.Code != tt.expectedCode {
						t.Errorf("expected error code %s, got %s", tt.expectedCode, errResp.Error.Code)
					}
				}
			}
		})
	}
}

func TestMergePullRequestHandler_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	tests := []struct {
		name           string
		requestBody    MergePullRequestRequest
		expectedStatus int
	}{
		{
			name:           "invalid request body",
			requestBody:    MergePullRequestRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/pullRequest/merge", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			MergePullRequestHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestReassignPullRequestHandler_Integration(t *testing.T) {
	if os.Getenv("POSTGRES_HOST") == "" {
		t.Skip("Skipping integration test: POSTGRES_HOST not set")
	}

	tests := []struct {
		name           string
		requestBody    ReassignPullRequestRequest
		expectedStatus int
	}{
		{
			name:           "invalid request body",
			requestBody:    ReassignPullRequestRequest{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/pullRequest/reassign", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ReassignPullRequestHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Unit tests for helper functions
func TestRequestToDomainPR(t *testing.T) {
	req := CreatePullRequestRequest{
		PullRequestID:   "pr1",
		PullRequestName: "Test PR",
		AuthorID:        "user1",
	}

	pr := requestToDomainPR(req)

	if pr.ID != "pr1" {
		t.Errorf("expected ID pr1, got %s", pr.ID)
	}
	if pr.Name != "Test PR" {
		t.Errorf("expected Name Test PR, got %s", pr.Name)
	}
	if pr.AuthorID != "user1" {
		t.Errorf("expected AuthorID user1, got %s", pr.AuthorID)
	}
	if pr.Status != domain.Open {
		t.Errorf("expected Status Open, got %v", pr.Status)
	}
}

func TestDomainPRToResponse(t *testing.T) {
	now := time.Now()
	pr := domain.PullRequest{
		PullRequestShort: domain.PullRequestShort{
			ID:       "pr1",
			Name:     "Test PR",
			AuthorID: "user1",
			Status:   domain.Open,
		},
		AssignedReviewers: "[user2, user3]",
		CreatedAt:         &now,
		MergedAt:          nil,
	}

	resp := domainPRToResponse(pr)

	if resp.ID != "pr1" {
		t.Errorf("expected ID pr1, got %s", resp.ID)
	}
	if len(resp.AssignedReviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(resp.AssignedReviewers))
	}
	if resp.AssignedReviewers[0] != "user2" && resp.AssignedReviewers[0] != "user3" {
		t.Errorf("unexpected reviewer: %s", resp.AssignedReviewers[0])
	}
}
