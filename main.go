package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/samber/do/v2"

	_ "github.com/mattn/go-sqlite3"
)

const BASE_URL = "https://asia-southeast2.cloud.twistlock.com/indonesia-751472959/api/v34.03"

func main() {
	injector := startProgram()
	defer func() {
		// Cleanup database connection
		db := do.MustInvoke[*sql.DB](injector)
		db.Close()
	}()

	service := do.MustInvoke[*Service](injector)

	mux := http.NewServeMux()

	// container endpoints
	mux.HandleFunc("/profile/container", fetchProfile(service))
	mux.HandleFunc("/policy/container", fetchPolicies(service))

	// host endpoints
	mux.HandleFunc("/profile/host", fetchHostPolicies(service))

	mux.HandleFunc("/verdict/send", sendVerdict(service))
	mux.HandleFunc("/verdict/update", updateVerdict(service))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Server starting on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /profile/container - Fetch and save container profiles")
	fmt.Println("  GET  /policy/container - Fetch and save runtime container policies")
	fmt.Println("  GET  /profile/host - Fetch and save runtime host policies")
	fmt.Println("  GET  /verdict/send - Send verdict email with CSV")
	fmt.Println("  POST /verdict/update - Update verdicts from CSV file")
	fmt.Println("  GET  /health - Health check")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(fmt.Errorf("Failed to start server: %v", err))
	}
}
