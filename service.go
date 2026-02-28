package main

import (
	"fmt"
	"os"
	"strings"
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

func (s *Service) FetchAndSaveHostPolicies() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	policy, err := getAllRuntimeHostPolicies(token)
	if err != nil {
		return fmt.Errorf("failed to get host policies: %v", err)
	}

	// Save policies to database
	err = s.Repo.SaveHostRules(policy)
	if err != nil {
		return fmt.Errorf("failed to save host policies: %v", err)
	}

	fmt.Printf("Successfully saved data from host policy %s with %d rules to database\n", policy.ID, len(policy.Rules))
	return nil
}

func (s *Service) FetchAndSaveHostProfiles() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	profiles, err := getRuntimeHostProfile(token)
	if err != nil {
		return fmt.Errorf("failed to get host profiles: %v", err)
	}

	// Transform profiles into records with business logic
	var records []HostProfileRecord
	for _, profile := range profiles {
		// Filter collections to exclude "All"
		collections := []string{}
		for _, col := range profile.Collections {
			if strings.ToLower(col) != "all" {
				collections = append(collections, col)
			}
		}

		// Skip if no valid collections
		if len(collections) == 0 {
			continue
		}

		for _, collection := range collections {
			// Process SSH events
			for _, ssh := range profile.SSHEvents {
				if ssh.Command != "" {
					value := fmt.Sprintf("command=%s,user=%s,ip=%d", ssh.Command, ssh.User, ssh.IP)
					records = append(records, HostProfileRecord{
						HostID:        profile.ID,
						CollectionName: collection,
						Key:           "ssh_event",
						Value:         value,
					})
				}
			}

			// Process app processes
			for _, app := range profile.Apps {
				// Save startup process
				if app.StartupProcess != nil && app.StartupProcess.Path != "" {
					records = append(records, HostProfileRecord{
						HostID:        profile.ID,
						CollectionName: collection,
						Key:           "process",
						Value:         app.StartupProcess.Path,
					})
				}

				// Save app processes
				for _, proc := range app.Processes {
					if proc.Path != "" {
						records = append(records, HostProfileRecord{
							HostID:        profile.ID,
							CollectionName: collection,
							Key:           "process",
							Value:         proc.Path,
						})
					}
				}

				// Save listening ports
				for _, lp := range app.ListeningPorts {
					if lp.Port > 0 {
						records = append(records, HostProfileRecord{
							HostID:        profile.ID,
							CollectionName: collection,
							Key:           "listening_port",
							Value:         fmt.Sprintf("%d", lp.Port),
						})
					}
				}

				// Save outgoing ports
				for _, op := range app.OutgoingPorts {
					if op.Port > 0 {
						records = append(records, HostProfileRecord{
							HostID:        profile.ID,
							CollectionName: collection,
							Key:           "outgoing_port",
							Value:         fmt.Sprintf("%d", op.Port),
						})
					}
				}
			}
		}
	}

	// Save records to database
	err = s.Repo.SaveHostProfileRecords(records)
	if err != nil {
		return fmt.Errorf("failed to save host profiles: %v", err)
	}

	fmt.Printf("Successfully saved data from %d host profiles to database\n", len(profiles))
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
