package main

import (
	"fmt"
	"log"

	emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
	// Replace with your actual API key
	apiKey := "your_api_key_here"

	// Initialize client
	client := emaillistchecker.NewClient(apiKey)

	// Verify an email
	fmt.Println("Verifying email...")
	result, err := client.Verify("test@example.com", nil, true)
	if err != nil {
		log.Fatal(err)
	}

	// Display results
	fmt.Println("\n=== Verification Result ===")
	fmt.Printf("Email: %s\n", result.Email)
	fmt.Printf("Result: %s\n", result.Result)
	fmt.Printf("Reason: %s\n", result.Reason)
	fmt.Printf("Score: %.2f\n", result.Score)

	fmt.Println("\n=== Email Details ===")
	fmt.Printf("Disposable: %t\n", result.Disposable)
	fmt.Printf("Role-based: %t\n", result.Role)
	fmt.Printf("Free provider: %t\n", result.Free)
	fmt.Printf("SMTP Provider: %s\n", result.SMTPProvider)
	fmt.Printf("Domain: %s\n", result.Domain)

	if len(result.MXRecords) > 0 {
		fmt.Println("\nMX Records:")
		for _, mx := range result.MXRecords {
			fmt.Printf("  - %s\n", mx)
		}
	}
}
