package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
)

func AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrorCodeNotFound, "invalid request body")
		return
	}

	existingTeam, err := application.GetTeam(r.Context(), req.TeamName)
	if err == nil && existingTeam.TeamName != "" {
		writeError(w, http.StatusBadRequest, ErrorCodeTeamExists, "team_name already exists")
		return
	}

	team := requestToDomainTeam(req)
	result, err := application.AddTeam(r.Context(), team)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "user with id") {
				writeError(w, http.StatusBadRequest, ErrorCodeNotFound, err.Error())
				return
			}
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

	team, err := application.GetTeam(r.Context(), teamName)
	if err != nil || team.TeamName == "" {
		writeError(w, http.StatusNotFound, ErrorCodeNotFound, "resource not found")
		return
	}

	writeJSON(w, http.StatusOK, team)
}
