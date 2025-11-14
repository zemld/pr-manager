package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func SetUserActiveHandler(w http.ResponseWriter, r *http.Request) {
	var req SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	user := domain.User{
		UserID:   req.UserID,
		IsActive: req.IsActive,
	}

	result, err := application.UpdateUserStatus(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	writeJSON(w, http.StatusOK, UserWrapperResponse{
		User: result,
	})
}

func GetUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "user_id parameter is required")
		return
	}

	prs, err := application.GetUserPullRequestsReviews(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	shortPRs := make([]PullRequestShortResponse, len(prs))
	for i, pr := range prs {
		shortPRs[i] = domainPRToShortResponse(pr)
	}

	writeJSON(w, http.StatusOK, UserPullRequestsResponse{
		UserID:       userID,
		PullRequests: shortPRs,
	})
}
