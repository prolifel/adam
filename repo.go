package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type Repo struct {
	DB *sql.DB
}

func (r *Repo) SaveProfiles(profiles []ContainerProfile) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO container_profiles (collection_name, key, value, verdict, updated_at)
		VALUES (?, ?, ?, 'not_yet', CURRENT_TIMESTAMP)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

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
			// Save DNS queries
			for _, dns := range profile.Network.Behavioral.DNSQueries {
				if dns.DomainName != "" {
					_, err := stmt.Exec(collection, "dns_queries", dns.DomainName)
					if err != nil {
						return err
					}
				}
			}

			// Save listening ports
			for _, lp := range profile.Network.Behavioral.ListeningPorts {
				for _, port := range lp.PortsData.Ports {
					value := fmt.Sprintf("%d", port.Port)
					_, err := stmt.Exec(collection, "listening_port", value)
					if err != nil {
						return err
					}
				}
			}

			// Save outbound ports
			for _, port := range profile.Network.Behavioral.OutboundPorts.Ports {
				value := fmt.Sprintf("%d", port.Port)
				_, err := stmt.Exec(collection, "outbound_port", value)
				if err != nil {
					return err
				}
			}

			// Save behavioral processes
			for _, proc := range profile.Processes.Behavioral {
				if proc.Path != "" {
					_, err := stmt.Exec(collection, "process", proc.Path)
					if err != nil {
						return err
					}
				}
			}

			// Save static processes
			for _, proc := range profile.Processes.Static {
				if proc.Path != "" {
					_, err := stmt.Exec(collection, "process", proc.Path)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return tx.Commit()
}

func (r *Repo) ExportNotYetVerdict() (string, error) {
	// Query records with "not_yet" verdict
	rows, err := r.DB.Query(`
		SELECT id, collection_name, key, value, verdict, COALESCE(remarks, '') as remarks
		FROM container_profiles
		WHERE verdict = 'not_yet'
		ORDER BY collection_name, key, value
	`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Create CSV file with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("container_profiles_not_yet_%s.csv", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"id", "collection_name", "key", "value", "verdict", "remarks"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write data rows
	count := 0
	for rows.Next() {
		var id int
		var collectionName, key, value, verdict, remarks string

		if err := rows.Scan(&id, &collectionName, &key, &value, &verdict, &remarks); err != nil {
			return "", err
		}

		row := []string{
			fmt.Sprintf("%d", id),
			collectionName,
			key,
			value,
			verdict,
			remarks,
		}

		if err := writer.Write(row); err != nil {
			return "", err
		}
		count++
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	fmt.Printf("Exported %d records to %s\n", count, filename)
	return filename, nil
}

func (r *Repo) UpdateVerdicts(records []CapabilitiesCSVHeader) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		UPDATE container_profiles 
		SET verdict = ?, remarks = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	updatedCount := 0
	for _, record := range records {
		// Validate verdict value
		if record.Verdict != "not_yet" && record.Verdict != "legitimate" && record.Verdict != "not_legitimate" {
			return 0, fmt.Errorf("invalid verdict value '%s' for ID %s. Must be: not_yet, legitimate, or not_legitimate", record.Verdict, record.ID)
		}

		result, err := stmt.Exec(record.Verdict, record.Remarks, record.ID)
		if err != nil {
			return 0, fmt.Errorf("failed to update record ID %s: %v", record.ID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}

		if rowsAffected > 0 {
			updatedCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	fmt.Printf("Updated %d records in database\n", updatedCount)
	return updatedCount, nil
}

func (r *Repo) SaveRules(policy ContainerPolicy) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO container_rules (collection_name, rule)
		VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Track which collections we've already saved to avoid duplicates
	savedCollections := make(map[string]bool)

	// If no rules, nothing to save
	if len(policy.Rules) == 0 {
		return tx.Commit()
	}

	// Iterate over rules and extract collection names
	for _, rule := range policy.Rules {
		// If no collections defined, use rule name as collection name
		if len(rule.Collections) == 0 {
			if rule.Name != "" && !savedCollections[rule.Name] {
				ruleJSON, err := json.Marshal(rule)
				if err != nil {
					return err
				}
				_, err = stmt.Exec(rule.Name, string(ruleJSON))
				if err != nil {
					return err
				}
				savedCollections[rule.Name] = true
			}
			continue
		}

		// Save each collection for this rule
		for _, collection := range rule.Collections {
			collectionName := collection.Name
			if collectionName == "" || collectionName == "Container - Alert on All Container" || collectionName == "All" {
				continue
			}
			if savedCollections[collectionName] {
				continue
			}

			ruleJSON, err := json.Marshal(rule)
			if err != nil {
				return err
			}
			_, err = stmt.Exec(collectionName, string(ruleJSON))
			if err != nil {
				return err
			}
			savedCollections[collectionName] = true
		}
	}

	return tx.Commit()
}

// GetRuleByCollection retrieves a single rule from container_rules by collection name
func (r *Repo) GetRuleByCollection(collectionName string) (ContainerRule, error) {
	var ruleJSON string
	err := r.DB.QueryRow(`
		SELECT rule FROM container_rules WHERE collection_name = ?
	`, collectionName).Scan(&ruleJSON)

	if err == sql.ErrNoRows {
		return ContainerRule{}, fmt.Errorf("no rule found for collection: %s", collectionName)
	}
	if err != nil {
		return ContainerRule{}, err
	}

	var rule ContainerRule
	if err := json.Unmarshal([]byte(ruleJSON), &rule); err != nil {
		return ContainerRule{}, fmt.Errorf("failed to unmarshal rule: %v", err)
	}

	return rule, nil
}

// UpdateRuleWithVerdict updates a rule in container_rules with a new verdict value
func (r *Repo) UpdateRuleWithVerdict(collectionName string, key string, value string) error {
	// Get existing rule
	rule, err := r.GetRuleByCollection(collectionName)
	if err != nil {
		// If rule doesn't exist, create a new one
		if strings.Contains(err.Error(), "no rule found") {
			rule = ContainerRule{
				Name: collectionName,
				DNS: DNSRule{
					DomainList: DomainList{
						Allowed: []string{},
						Denied:  []string{},
					},
				},
				Processes: ProcessRule{
					AllowedList: []string{},
					DeniedList:  DeniedList{Paths: []string{}},
				},
			}
		} else {
			return err
		}
	}

	// Update the rule based on the key type
	switch key {
	case "dns_queries":
		// Check if domain already exists in allowed list
		if !slices.Contains(rule.DNS.DomainList.Allowed, value) {
			rule.DNS.DomainList.Allowed = append(rule.DNS.DomainList.Allowed, value)
			fmt.Printf("Added DNS allowed: %s to collection %s\n", value, collectionName)
		}
	case "processes", "process":
		// Check if process already exists in allowed list
		if !slices.Contains(rule.Processes.AllowedList, value) {
			rule.Processes.AllowedList = append(rule.Processes.AllowedList, value)
			fmt.Printf("Added process allowed: %s to collection %s\n", value, collectionName)
		}
	default:
		return fmt.Errorf("unknown key type: %s", key)
	}

	// Marshal the updated rule back to JSON
	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return fmt.Errorf("failed to marshal rule: %v", err)
	}

	// Update in database
	_, err = r.DB.Exec(`
		INSERT OR REPLACE INTO container_rules (collection_name, rule)
		VALUES (?, ?)
	`, collectionName, string(ruleJSON))

	if err != nil {
		return fmt.Errorf("failed to update rule: %v", err)
	}

	return nil
}

// GetAllRules retrieves all rules from container_rules table
func (r *Repo) GetAllRules() ([]ContainerRule, error) {
	rows, err := r.DB.Query(`
		SELECT rule FROM container_rules
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []ContainerRule
	for rows.Next() {
		var ruleJSON string
		if err := rows.Scan(&ruleJSON); err != nil {
			return nil, err
		}

		var rule ContainerRule
		if err := json.Unmarshal([]byte(ruleJSON), &rule); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rule: %v", err)
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}
