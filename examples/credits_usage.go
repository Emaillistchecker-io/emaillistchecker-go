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

	// Get credit balance
	fmt.Println("=== Credit Balance ===")
	credits, err := client.GetCredits()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Available credits: %v\n", credits["balance"])
	fmt.Printf("Used this month: %v\n", credits["used_this_month"])
	fmt.Printf("Current plan: %v\n\n", credits["plan"])

	// Get usage statistics
	fmt.Println("=== Usage Statistics ===")
	usage, err := client.GetUsage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total API requests: %v\n", usage["total_requests"])
	fmt.Printf("Successful requests: %v\n", usage["successful_requests"])
	fmt.Printf("Failed requests: %v\n", usage["failed_requests"])

	// Calculate success rate
	if totalRequests, ok := usage["total_requests"].(float64); ok && totalRequests > 0 {
		if successfulRequests, ok := usage["successful_requests"].(float64); ok {
			successRate := (successfulRequests / totalRequests) * 100
			fmt.Printf("Success rate: %.2f%%\n", successRate)
		}
	}
}
