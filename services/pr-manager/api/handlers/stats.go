package handlers

import (
	"net/http"

	"github.com/zemld/pr-manager/pr-manager/internal/application"
)

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := application.GetStats(r.Context())
	writeJSON(w, http.StatusOK, stats)
}
