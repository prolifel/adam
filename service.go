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
			// NOTE: sekarang masih belum perlu
			// for _, ssh := range profile.SSHEvents {
			// 	if ssh.Command != "" {
			// 		value := fmt.Sprintf("command=%s,user=%s,ip=%d", ssh.Command, ssh.User, ssh.IP)
			// 		records = append(records, HostProfileRecord{
			// 			HostID:        profile.ID,
			// 			CollectionName: collection,
			// 			Key:           "ssh_event",
			// 			Value:         value,
			// 		})
			// 	}
			// }

			// Process app processes
			for _, app := range profile.Apps {
				// Save startup process
				if app.StartupProcess != nil && app.StartupProcess.Path != "" {
					records = append(records, HostProfileRecord{
						HostID:         profile.ID,
						CollectionName: collection,
						Key:            "process",
						Value:          app.StartupProcess.Path,
					})
				}

				// Save app processes
				for _, proc := range app.Processes {
					if proc.Path != "" {
						records = append(records, HostProfileRecord{
							HostID:         profile.ID,
							CollectionName: collection,
							Key:            "process",
							Value:          proc.Path,
						})
					}
				}

				// Save listening ports
				for _, lp := range app.ListeningPorts {
					if lp.Port > 0 {
						records = append(records, HostProfileRecord{
							HostID:         profile.ID,
							CollectionName: collection,
							Key:            "listening_port",
							Value:          fmt.Sprintf("%d", lp.Port),
						})
					}
				}

				// Save outgoing ports
				for _, op := range app.OutgoingPorts {
					if op.Port > 0 {
						records = append(records, HostProfileRecord{
							HostID:         profile.ID,
							CollectionName: collection,
							Key:            "outgoing_port",
							Value:          fmt.Sprintf("%d", op.Port),
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

func (s *Service) FetchAndSaveAppEmbeddedProfiles() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	profiles, err := getAppEmbeddedProfile(token)
	if err != nil {
		return fmt.Errorf("failed to get app-embedded profiles: %v", err)
	}

	// Transform profiles into records
	var records []AppEmbeddedProfileRecord
	for _, profile := range profiles {
		// If no collections, skip
		if len(profile.Collections) == 0 {
			continue
		}

		// Save profile_id for each collection
		for _, collection := range profile.Collections {
			// Skip "All" collection
			if collection == "All" {
				continue
			}

			// Save profile_id as key
			records = append(records, AppEmbeddedProfileRecord{
				ProfileID:     profile.ID,
				AppID:         profile.AppID,
				CollectionName: collection,
				Key:           "profile_id",
				Value:         profile.ID,
			})

			// Save app_id if present
			if profile.AppID != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "app_id",
					Value:         profile.AppID,
				})
			}

			// Save cluster if present
			if profile.Cluster != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "cluster",
					Value:         profile.Cluster,
				})
			}

			// Save container if present
			if profile.Container != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "container",
					Value:         profile.Container,
				})
			}

			// Save image if present
			if profile.Image != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "image",
					Value:         profile.Image,
				})
			}

			// Save imageID if present
			if profile.ImageID != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "image_id",
					Value:         profile.ImageID,
				})
			}

			// Save startTime if present
			if profile.StartTime != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "start_time",
					Value:         profile.StartTime,
				})
			}

			// Save clusterType if present
			if profile.ClusterType != "" {
				records = append(records, AppEmbeddedProfileRecord{
					ProfileID:     profile.ID,
					AppID:         profile.AppID,
					CollectionName: collection,
					Key:           "cluster_type",
					Value:         profile.ClusterType,
				})
			}
		}
	}

	// Save records to database
	err = s.Repo.SaveAppEmbeddedProfiles(records)
	if err != nil {
		return fmt.Errorf("failed to save app-embedded profiles: %v", err)
	}

	fmt.Printf("Successfully saved data from %d app-embedded profiles to database\n", len(profiles))
	return nil
}

func (s *Service) FetchAndSaveAppEmbeddedPolicies() error {
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	policy, err := getAppEmbeddedPolicy(token)
	if err != nil {
		return fmt.Errorf("failed to get app-embedded policies: %v", err)
	}

	// Save policies to database
	err = s.Repo.SaveAppEmbeddedRules(policy)
	if err != nil {
		return fmt.Errorf("failed to save app-embedded policies: %v", err)
	}

	fmt.Printf("Successfully saved data from app-embedded policy %s with %d rules to database\n", policy.ID, len(policy.Rules))
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

// GenerateWeeklyAlertReport generates the weekly CSPM alert report
func (s *Service) GenerateWeeklyAlertReport() (int, int, error) {
	// Login to Prisma Cloud
	token, err := login(s.Cfg.AccessKeyId, s.Cfg.SecretAccessKey)
	if err != nil {
		return 0, 0, fmt.Errorf("login failed: %v", err)
	}

	// Fetch AWS alerts
	awsAlerts, err := getCSPMAlerts(token, s.Cfg.ComplianceStandard, "AWS", true)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch AWS alerts: %v", err)
	}

	// Fetch GCP alerts
	gcpAlerts, err := getCSPMAlerts(token, s.Cfg.ComplianceStandard, "GCP", true)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch GCP alerts: %v", err)
	}

	// Combine alerts for CSV generation
	allAlerts := append(awsAlerts, gcpAlerts...)

	// Generate CSV files
	awsFile, gcpFile, err := generateAWSAndGCPCSVs(allAlerts)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to generate CSV files: %v", err)
	}

	// Send email with attachments
	err = sendAlertEmailWithCSVs(s.Cfg, awsFile, gcpFile, s.Cfg.ComplianceStandard, len(awsAlerts), len(gcpAlerts))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send email: %v", err)
	}

	// Cleanup CSV files after sending
	if awsFile != "" {
		if err := os.Remove(awsFile); err != nil {
			fmt.Printf("Warning: failed to delete AWS CSV file: %v\n", err)
		}
	}
	if gcpFile != "" {
		if err := os.Remove(gcpFile); err != nil {
			fmt.Printf("Warning: failed to delete GCP CSV file: %v\n", err)
		}
	}

	fmt.Printf("Weekly alert report completed: AWS=%d, GCP=%d\n", len(awsAlerts), len(gcpAlerts))
	return len(awsAlerts), len(gcpAlerts), nil
}
