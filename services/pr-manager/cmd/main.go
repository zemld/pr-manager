package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zemld/pr-manager/pr-manager/api/handlers"
	"github.com/zemld/pr-manager/pr-manager/internal/application"
)

func main() {
	ctx := context.Background()
	if err := application.InitializeDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", handlers.AddTeamHandler)
	mux.HandleFunc("GET /team/get", handlers.GetTeamHandler)
	mux.HandleFunc("DELETE /team/delete", handlers.DeleteTeamHandler)

	mux.HandleFunc("POST /users/setIsActive", handlers.SetUserActiveHandler)
	mux.HandleFunc("GET /users/getReview", handlers.GetUserReviewsHandler)

	mux.HandleFunc("POST /pullRequest/create", handlers.CreatePullRequestHandler)
	mux.HandleFunc("POST /pullRequest/merge", handlers.MergePullRequestHandler)
	mux.HandleFunc("POST /pullRequest/reassign", handlers.ReassignPullRequestHandler)

	mux.HandleFunc("GET /stats/get", handlers.GetStatsHandler)

	port := "8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
