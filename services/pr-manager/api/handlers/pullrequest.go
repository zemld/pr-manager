package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
)

func CreatePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	pr := requestToDomainPR(req)
	result, err := application.CreatePullRequest(r.Context(), pr)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			writeError(w, http.StatusConflict, ErrorCodePRExists, "PR id already exists")
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no possible assigners") {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, PullRequestWrapperResponse{
		PR: domainPRToResponse(result),
	})
}

func MergePullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	pr := requestToDomainPRForMerge(req)
	result, err := application.MergePullRequest(r.Context(), pr)
	if err != nil {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	writeJSON(w, http.StatusOK, PullRequestWrapperResponse{
		PR: domainPRToResponse(result),
	})
}

func ReassignPullRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req ReassignPullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	result, newReviewer, err := application.ReassignPullRequest(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		if strings.Contains(err.Error(), "already merged") || strings.Contains(err.Error(), "merged") {
			writeError(w, http.StatusConflict, ErrorCodePRMerged, "cannot reassign on merged PR")
			return
		}
		if strings.Contains(err.Error(), "not assigned") {
			writeError(w, http.StatusConflict, ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
			return
		}
		if strings.Contains(err.Error(), "no active replacement") || strings.Contains(err.Error(), "no possible") {
			writeError(w, http.StatusConflict, ErrorCodeNoCandidate, "no active replacement candidate in team")
			return
		}
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	prResponse := domainPRToResponse(result)

	writeJSON(w, http.StatusOK, ReassignResponse{
		PR:         prResponse,
		ReplacedBy: newReviewer,
	})
}
