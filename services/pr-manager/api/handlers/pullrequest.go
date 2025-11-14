package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
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
		if errors.Is(err, domain.ErrPRExists) {
			writeError(w, http.StatusConflict, ErrorCodePRExists, "PR id already exists")
			return
		}
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrNoPossibleAssigners) {
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
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
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
		if errors.Is(err, domain.ErrPRMerged) {
			writeError(w, http.StatusConflict, ErrorCodePRMerged, "cannot reassign on merged PR")
			return
		}
		if errors.Is(err, domain.ErrNotAssigned) {
			writeError(w, http.StatusConflict, ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
			return
		}
		if errors.Is(err, domain.ErrNoCandidate) {
			writeError(w, http.StatusConflict, ErrorCodeNoCandidate, "no active replacement candidate in team")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}

	prResponse := domainPRToResponse(result)

	writeJSON(w, http.StatusOK, ReassignResponse{
		PR:         prResponse,
		ReplacedBy: newReviewer,
	})
}
