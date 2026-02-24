package main

import (
	"fmt"
	"os"
)

type Service struct {
	Repo *Repo
	Cfg  Config
}

func (s *Service) SendVerdict() error {
	// Export CSV
	csvFilename, err := s.Repo.ExportNotYetVerdict()
	if err != nil {
		return fmt.Errorf("failed to export CSV: %v", err)
	}

	// Send email
	if err := sendEmailWithCSV(s.Cfg, csvFilename); err != nil {
		return err
	}

	// Delete the CSV file after sending
	if err := os.Remove(csvFilename); err != nil {
		fmt.Printf("Warning: failed to delete CSV file: %v\n", err)
	}

	return nil
}

func (s *Service) FetchAndSaveProfiles() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	profiles, err := getRuntimeContainerProfile(token)
	if err != nil {
		return fmt.Errorf("failed to get profiles: %v", err)
	}

	// Save profiles to database
	err = s.Repo.SaveProfiles(profiles)
	if err != nil {
		return fmt.Errorf("failed to save profiles: %v", err)
	}

	fmt.Printf("Successfully saved data from %d profiles to database\n", len(profiles))
	return nil
}
