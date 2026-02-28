package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// detectDelimiter auto-detects the CSV delimiter by counting occurrences
func detectDelimiter(data []byte) rune {
	comma := bytes.Count(data, []byte(","))
	semicolon := bytes.Count(data, []byte(";"))
	tab := bytes.Count(data, []byte("\t"))

	if semicolon > comma && semicolon > tab {
		return ';'
	}
	if tab > comma && tab > semicolon {
		return '\t'
	}
	return ',' // default
}

// parseCSVWithAutoDetect parses a CSV file with auto-detected delimiter
func parseCSVWithAutoDetect(file io.Reader) ([]CapabilitiesCSVHeader, error) {
	// Read entire file into buffer to detect delimiter
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	delimiter := detectDelimiter(data)
	fmt.Printf("Detected CSV delimiter: %c\n", delimiter)

	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = delimiter

	// Read header first (and discard)
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var records []CapabilitiesCSVHeader
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// record format: id, collection_name, key, value, verdict, remarks
		if len(record) < 5 {
			continue
		}

		records = append(records, CapabilitiesCSVHeader{
			ID:             strings.TrimSpace(record[0]),
			CollectionName: strings.TrimSpace(record[1]),
			Key:            strings.TrimSpace(record[2]),
			Value:          strings.TrimSpace(record[3]),
			Verdict:        strings.TrimSpace(record[4]),
			Remarks:        strings.TrimSpace(record[5]),
		})
	}

	return records, nil
}

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

		capabilities, err := parseCSVWithAutoDetect(file)
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

		// Push legitimate verdicts to Prisma Cloud
		pcPushedCount, err := service.PushVerdictToPrismaCloud(capabilities)
		if err != nil {
			// Log error but don't fail the request - local DB update succeeded
			fmt.Printf("Warning: Failed to push to Prisma Cloud: %v\n", err)
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: fmt.Sprintf("Successfully updated %d records", updatedCount),
			Data: map[string]int{
				"updated_count":  updatedCount,
				"total_records":  len(capabilities),
				"pc_pushed_count": pcPushedCount,
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

func fetchPolicies(service *Service) http.HandlerFunc {
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

		err = service.FetchAndSavePolicies()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch policies: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: "Policies fetched and saved successfully!",
		}

		res, err := json.Marshal(resp)
		if err != nil {
			return
		}

		w.Write(res)
	}
}

func fetchHostPolicies(service *Service) http.HandlerFunc {
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

		err = service.FetchAndSaveHostPolicies()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch host policies: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp := Response{
			Message: "Host policies fetched and saved successfully!",
		}

		res, err := json.Marshal(resp)
		if err != nil {
			return
		}

		w.Write(res)
	}
}
