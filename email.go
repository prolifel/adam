package main

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func sendEmailWithCSV(cfg Config, csvFilename string) error {
	// Parse recipient emails
	recipients := strings.Split(cfg.EmailTo, ",")
	for i, email := range recipients {
		recipients[i] = strings.TrimSpace(email)
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailFrom)
	m.SetHeader("To", recipients...)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	m.SetHeader("Subject", fmt.Sprintf("Container Profiles Review - %s", timestamp))

	body := fmt.Sprintf(`Hello,

Please find attached the container profiles that require review.

This CSV file contains all entries with verdict status "not_yet".

Timestamp: %s
File: %s

Best regards,
Adam`, timestamp, csvFilename)

	m.SetBody("text/plain", body)
	m.Attach(csvFilename)

	// Send email
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("Email sent successfully to: %s\n", strings.Join(recipients, ", "))
	return nil
}
