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

func (s *Service) FetchAndSavePolicies() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	policy, err := getAllRuntimeContainerPolicies(token)
	if err != nil {
		return fmt.Errorf("failed to get policies: %v", err)
	}

	// Save policies to database
	err = s.Repo.SaveRules(policy)
	if err != nil {
		return fmt.Errorf("failed to save policies: %v", err)
	}

	fmt.Printf("Successfully saved data from policy %s with %d rules to database\n", policy.ID, len(policy.Rules))
	return nil
}

func (s *Service) PushVerdictToPrismaCloud(verdicts []CapabilitiesCSVHeader) (int, error) {
	// Filter only legitimate verdicts
	var legitimateVerdicts []CapabilitiesCSVHeader
	for _, v := range verdicts {
		if v.Verdict == "legitimate" {
			legitimateVerdicts = append(legitimateVerdicts, v)
		}
	}

	if len(legitimateVerdicts) == 0 {
		fmt.Println("No legitimate verdicts to push to Prisma Cloud")
		return 0, nil
	}

	// Update rules in database for each verdict
	addedCount := 0
	for _, v := range legitimateVerdicts {
		err := s.Repo.UpdateRuleWithVerdict(v.CollectionName, v.Key, v.Value)
		if err != nil {
			return 0, fmt.Errorf("failed to update rule for collection %s: %v", v.CollectionName, err)
		}
		addedCount++
	}

	// Get all rules from the database
	rules, err := s.Repo.GetAllRules()
	if err != nil {
		return 0, fmt.Errorf("failed to get all rules: %v", err)
	}

	if len(rules) == 0 {
		fmt.Println("No rules found in database to push")
		return 0, nil
	}

	// Login to get token
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return 0, fmt.Errorf("login failed: %v", err)
	}

	// Build policy with all rules from database
	policy := ContainerPolicy{
		ID:    "containerRuntime",
		Rules: rules,
	}

	// Push updated policy back to Prisma Cloud
	err = updateRuntimeContainerPolicy(token, policy)
	if err != nil {
		return 0, fmt.Errorf("failed to update runtime container policy: %v", err)
	}

	fmt.Printf("Successfully pushed %d new rules to Prisma Cloud\n", addedCount)
	return addedCount, nil
}
