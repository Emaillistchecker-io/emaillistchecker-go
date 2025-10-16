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

	// Example 1: Find email by name and domain
	fmt.Println("=== Find Email by Name ===")
	result, err := client.FindEmail("John", "Doe", "example.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Email found: %v\n", result["email"])
	fmt.Printf("Confidence: %v%%\n", result["confidence"])
	fmt.Printf("Pattern: %v\n", result["pattern"])
	fmt.Printf("Verified: %v\n", result["verified"])

	if alternatives, ok := result["alternatives"].([]interface{}); ok && len(alternatives) > 0 {
		fmt.Println("\nAlternative patterns:")
		for _, alt := range alternatives {
			fmt.Printf("  - %v\n", alt)
		}
	}

	fmt.Println()

	// Example 2: Find emails by domain
	fmt.Println("=== Find Emails by Domain ===")
	domainResults, err := client.FindByDomain("example.com", 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Domain: %v\n", domainResults["domain"])
	fmt.Printf("Total found: %v\n", domainResults["total_found"])

	if patterns, ok := domainResults["patterns"].([]interface{}); ok && len(patterns) > 0 {
		fmt.Println("\nCommon email patterns:")
		for _, pattern := range patterns {
			fmt.Printf("  - %v\n", pattern)
		}
	}

	fmt.Println()

	// Example 3: Find emails by company
	fmt.Println("=== Find Emails by Company ===")
	companyResults, err := client.FindByCompany("Acme Corporation", 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Company: %v\n", companyResults["company"])
	fmt.Printf("Total found: %v\n", companyResults["total_found"])

	if domains, ok := companyResults["possible_domains"].([]interface{}); ok && len(domains) > 0 {
		fmt.Println("\nPossible domains:")
		for _, domain := range domains {
			fmt.Printf("  - %v\n", domain)
		}
	}
}
