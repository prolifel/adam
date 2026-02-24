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

	// Fetch profiles endpoint
	mux.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := service.FetchAndSaveProfiles()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch profiles: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profiles fetched and saved successfully"))
	})

	// Send verdict email endpoint
	mux.HandleFunc("/send-verdict", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := service.SendVerdict()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send verdict email: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verdict email sent successfully"))
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Server starting on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  POST /fetch - Fetch and save container profiles")
	fmt.Println("  POST /send-verdict - Send verdict email with CSV")
	fmt.Println("  GET  /health - Health check")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(fmt.Errorf("Failed to start server: %v", err))
	}
}
