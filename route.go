package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
