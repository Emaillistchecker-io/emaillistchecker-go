package main

import (
	"fmt"
	"log"
	"time"

	emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
	// Replace with your actual API key
	apiKey := "your_api_key_here"

	// Initialize client
	client := emaillistchecker.NewClient(apiKey)

	// List of emails to verify
	emails := []string{
		"user1@example.com",
		"user2@example.com",
		"user3@example.com",
		"invalid@invalid-domain-xyz.com",
		"test@gmail.com",
	}

	fmt.Printf("Submitting batch of %d emails...\n", len(emails))

	// Submit batch
	batch, err := client.VerifyBatch(emails, "My Test Batch", "", true)
	if err != nil {
		log.Fatal(err)
	}

	batchID := batch.ID
	fmt.Println("Batch submitted successfully!")
	fmt.Printf("Batch ID: %d\n", batchID)
	fmt.Printf("Status: %s\n", batch.Status)
	fmt.Printf("Total emails: %d\n\n", batch.TotalEmails)

	// Monitor progress
	fmt.Println("Monitoring progress...")
	previousProgress := 0

	for {
		status, err := client.GetBatchStatus(batchID)
		if err != nil {
			log.Fatal(err)
		}

		if status.Progress != previousProgress {
			fmt.Printf("Progress: %d%% (%d/%d processed)\n",
				status.Progress, status.ProcessedEmails, status.TotalEmails)
			previousProgress = status.Progress
		}

		if status.Status == "completed" {
			fmt.Println("\nBatch verification completed!\n")
			break
		} else if status.Status == "failed" {
			fmt.Println("\nBatch verification failed!")
			return
		}

		time.Sleep(2 * time.Second) // Wait 2 seconds before checking again
	}

	// Get final statistics
	finalStatus, err := client.GetBatchStatus(batchID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Final Statistics ===")
	fmt.Printf("Total: %d\n", finalStatus.TotalEmails)
	fmt.Printf("Valid: %d\n", finalStatus.ValidEmails)
	fmt.Printf("Invalid: %d\n", finalStatus.InvalidEmails)
	fmt.Printf("Unknown: %d\n\n", finalStatus.UnknownEmails)

	// Download results
	fmt.Println("Downloading results...")
	results, err := client.GetBatchResults(batchID, "json", "all")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== Results ===")
	fmt.Printf("%+v\n", results)
}
