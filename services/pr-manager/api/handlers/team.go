package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
	"github.com/zemld/pr-manager/pr-manager/internal/domain"
)

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	teamName := req.TeamName
	existingTeam, err := application.GetTeam(r.Context(), &teamName)
	if err == nil && existingTeam.TeamName != "" {
		writeError(w, http.StatusBadRequest, ErrorCodeTeamExists, "team_name already exists")
		return
	}

	team := requestToDomainTeam(req)
	result, err := application.AddTeam(r.Context(), team)
	if err != nil {
		if errors.Is(err, domain.ErrUserInAnotherTeam) {
			writeError(w, http.StatusBadRequest, ErrorCodeNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrTeamExists) {
			writeError(w, http.StatusBadRequest, ErrorCodeTeamExists, "team_name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, TeamWrapperResponse{
		Team: result,
	})
}

func GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "name parameter is required")
		return
	}

	team, err := application.GetTeam(r.Context(), &teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) || errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}
	if team.TeamName == "" {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	writeJSON(w, http.StatusOK, team)
}

func DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "name parameter is required")
		return
	}

	existingTeam, err := application.GetTeam(r.Context(), &teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) || errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}
	if existingTeam.TeamName == "" {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	err = application.DeleteTeam(r.Context(), teamName)
	if err != nil {
		if errors.Is(err, domain.ErrTeamNotFound) || errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
			return
		}
		writeError(w, http.StatusInternalServerError, ErrorCodeNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
