package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nao1215/fileprep"
)

func fetchProfile(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := validateToken(r, service.Cfg.Token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = service.FetchAndSaveProfiles()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch profiles: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: "Profiles fetched and saved successfully!",
		}

		res, err := json.Marshal(resp)
		if err != nil {
			return
		}

		w.Write(res)
	}
}

func sendVerdict(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := validateToken(r, service.Cfg.Token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = service.SendVerdict()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send verdict email: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: "Verdict email sent successfully",
		}

		res, err := json.Marshal(resp)
		if err != nil {
			return
		}

		w.Write(res)
	}
}

func updateVerdict(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := validateToken(r, service.Cfg.Token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File not found in request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		processor := fileprep.NewProcessor(fileprep.FileTypeCSV)
		var capabilities []CapabilitiesCSVHeader

		_, _, err = processor.Process(file, &capabilities)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to process CSV: %v", err), http.StatusBadRequest)
			return
		}

		// Update verdicts in database
		updatedCount, err := service.Repo.UpdateVerdicts(capabilities)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update verdicts: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: fmt.Sprintf("Successfully updated %d records", updatedCount),
			Data: map[string]int{
				"updated_count": updatedCount,
				"total_records": len(capabilities),
			},
		}

		res, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}
