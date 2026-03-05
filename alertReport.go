package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// generateAlertCSV creates a CSV file from a list of CSPM alerts
func generateAlertCSV(alerts []CSPMAlert, filename string) error {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header (without Severity)
	header := []string{"ID", "Title", "Status", "Resource", "Policy", "CloudType", "AccountID", "Region", "CreatedTime", "Recommendation"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write data rows
	for _, alert := range alerts {
		row := []string{
			alert.ID,
			alert.Title,
			alert.Status,
			alert.Resource,
			alert.Policy,
			alert.CloudType,
			alert.AccountID,
			alert.Region,
			alert.CreatedTime,
			alert.Recommendation,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	fmt.Printf("Generated CSV file: %s with %d alerts\n", filename, len(alerts))
	return nil
}

// generateAWSAndGCPCSVs splits alerts by cloud type and generates separate CSV files
func generateAWSAndGCPCSVs(alerts []CSPMAlert) (awsFile, gcpFile string, err error) {
	// Get current date for filename
	date := time.Now().Format("2006-01-02")

	// Split alerts by cloud type
	var awsAlerts []CSPMAlert
	var gcpAlerts []CSPMAlert

	for _, alert := range alerts {
		if alert.CloudType == "AWS" {
			awsAlerts = append(awsAlerts, alert)
		} else if alert.CloudType == "GCP" {
			gcpAlerts = append(gcpAlerts, alert)
		}
	}

	// Generate AWS CSV
	if len(awsAlerts) > 0 {
		awsFilename := fmt.Sprintf("aws_alerts_%s.csv", date)
		if err := generateAlertCSV(awsAlerts, awsFilename); err != nil {
			return "", "", err
		}
		awsFile = awsFilename
	}

	// Generate GCP CSV
	if len(gcpAlerts) > 0 {
		gcpFilename := fmt.Sprintf("gcp_alerts_%s.csv", date)
		if err := generateAlertCSV(gcpAlerts, gcpFilename); err != nil {
			return "", "", err
		}
		gcpFile = gcpFilename
	}

	fmt.Printf("Generated CSV files - AWS: %d alerts, GCP: %d alerts\n", len(awsAlerts), len(gcpAlerts))
	return awsFile, gcpFile, nil
}
