package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
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
