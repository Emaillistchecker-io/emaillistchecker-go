# EmailListChecker Go SDK

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Official Go SDK for the [EmailListChecker](https://emaillistchecker.io) email verification API.

## Features

- **Email Verification** - Verify single or bulk email addresses
- **Email Finder** - Discover email addresses by name, domain, or company
- **Credit Management** - Check balance and usage
- **Batch Processing** - Async verification of large lists
- **Pure Go** - No external dependencies (uses standard library)
- **Type Safe** - Strongly typed structs
- **Error Handling** - Comprehensive error types

## Requirements

- Go 1.19 or higher

## Installation

```bash
go get github.com/Emaillistchecker-io/emaillistchecker-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    // Initialize client
    client := emaillistchecker.NewClient("your_api_key_here")

    // Verify an email
    result, err := client.Verify("test@example.com", nil, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Result: %s\n", result.Result)  // deliverable, undeliverable, risky, unknown
    fmt.Printf("Score: %.2f\n", result.Score)   // 0.0 to 1.0
}
```

## Get Your API Key

1. Sign up at [platform.emaillistchecker.io](https://platform.emaillistchecker.io/register)
2. Get your API key from the [API Dashboard](https://platform.emaillistchecker.io/api)
3. Start verifying!

## Usage Examples

### Single Email Verification

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    // Verify single email
    result, err := client.Verify("user@example.com", nil, true)
    if err != nil {
        log.Fatal(err)
    }

    switch result.Result {
    case "deliverable":
        fmt.Println("✓ Email is valid and deliverable")
    case "undeliverable":
        fmt.Println("✗ Email is invalid")
    case "risky":
        fmt.Println("⚠ Email is risky (catch-all, disposable, etc.)")
    default:
        fmt.Println("? Unable to determine")
    }

    // Check details
    fmt.Printf("Disposable: %t\n", result.Disposable)
    fmt.Printf("Role account: %t\n", result.Role)
    fmt.Printf("Free provider: %t\n", result.Free)
    fmt.Printf("SMTP provider: %s\n", result.SMTPProvider)
}
```

### Batch Email Verification

```go
package main

import (
    "fmt"
    "log"
    "time"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    // Submit batch for verification
    emails := []string{
        "user1@example.com",
        "user2@example.com",
        "user3@example.com",
    }

    batch, err := client.VerifyBatch(emails, "My Campaign List", "", true)
    if err != nil {
        log.Fatal(err)
    }

    batchID := batch.ID
    fmt.Printf("Batch ID: %d\n", batchID)
    fmt.Printf("Status: %s\n", batch.Status)

    // Check progress
    for {
        status, err := client.GetBatchStatus(batchID)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Progress: %d%%\n", status.Progress)

        if status.Status == "completed" {
            break
        }

        time.Sleep(5 * time.Second)  // Wait 5 seconds before checking again
    }

    // Download results
    results, err := client.GetBatchResults(batchID, "json", "all")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Results: %+v\n", results)
}
```

### Email Finder

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    // Find email by name and domain
    result, err := client.FindEmail("John", "Doe", "example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found: %s\n", result["email"])
    fmt.Printf("Confidence: %v%%\n", result["confidence"])
    fmt.Printf("Verified: %v\n", result["verified"])

    // Find all emails for a domain
    domainResults, err := client.FindByDomain("example.com", 50, 0)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Domain results: %+v\n", domainResults)

    // Find emails by company name
    companyResults, err := client.FindByCompany("Acme Corporation", 10)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Company results: %+v\n", companyResults)
}
```

### Credit Management

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    // Check credit balance
    credits, err := client.GetCredits()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Available credits: %v\n", credits["balance"])
    fmt.Printf("Used this month: %v\n", credits["used_this_month"])
    fmt.Printf("Current plan: %v\n", credits["plan"])

    // Get usage statistics
    usage, err := client.GetUsage()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total API calls: %v\n", usage["total_requests"])
    fmt.Printf("Successful: %v\n", usage["successful_requests"])
    fmt.Printf("Failed: %v\n", usage["failed_requests"])
}
```

### List Management

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    // Get all lists
    lists, err := client.GetLists()
    if err != nil {
        log.Fatal(err)
    }

    for _, list := range lists {
        listMap := list.(map[string]interface{})
        fmt.Printf("ID: %v\n", listMap["id"])
        fmt.Printf("Name: %v\n", listMap["name"])
        fmt.Printf("Status: %v\n", listMap["status"])
        fmt.Printf("Total emails: %v\n", listMap["total_emails"])
        fmt.Printf("Valid: %v\n", listMap["valid_emails"])
        fmt.Println("---")
    }

    // Delete a list
    err = client.DeleteList(123)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Error Handling

```go
package main

import (
    "fmt"
    "log"

    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
    client := emaillistchecker.NewClient("your_api_key")

    result, err := client.Verify("test@example.com", nil, true)
    if err != nil {
        switch e := err.(type) {
        case *emaillistchecker.AuthenticationError:
            fmt.Println("Invalid API key")
        case *emaillistchecker.InsufficientCreditsError:
            fmt.Println("Not enough credits")
        case *emaillistchecker.RateLimitError:
            fmt.Printf("Rate limit exceeded. Retry after %d seconds\n", e.RetryAfter)
        case *emaillistchecker.ValidationError:
            fmt.Printf("Validation error: %s\n", e.Message)
        case *emaillistchecker.Error:
            fmt.Printf("API error: %s (status: %d)\n", e.Message, e.StatusCode)
        default:
            log.Fatal(err)
        }
        return
    }

    fmt.Printf("Result: %s\n", result.Result)
}
```

## Configuration

### Custom Timeout

```go
import (
    "time"
    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

// Set custom timeout (default: 30 seconds)
client := emaillistchecker.NewClientWithConfig(
    "your_api_key",
    "https://platform.emaillistchecker.io/api/v1",
    60 * time.Second,  // 60 seconds timeout
)
```

### Custom Base URL

```go
import (
    "time"
    emaillistchecker "github.com/Emaillistchecker-io/emaillistchecker-go"
)

// Use custom API endpoint (for testing or private instances)
client := emaillistchecker.NewClientWithConfig(
    "your_api_key",
    "https://custom-api.example.com/api/v1",
    30 * time.Second,
)
```

## API Response Types

### Verification Result

```go
type VerifyResponse struct {
    Email        string   // Email address verified
    Result       string   // deliverable | undeliverable | risky | unknown
    Reason       string   // VALID | INVALID | ACCEPT_ALL | DISPOSABLE | etc.
    Disposable   bool     // Is temporary/disposable email
    Role         bool     // Is role-based (info@, support@, etc.)
    Free         bool     // Is free provider (gmail, yahoo, etc.)
    Score        float64  // Deliverability score (0.0 - 1.0)
    SMTPProvider string   // Email provider
    MXRecords    []string // List of MX records
    Domain       string   // Email domain
    SpamTrap     bool     // Is spam trap
    MXFound      bool     // MX records found
}
```

## Support

- **Documentation**: [platform.emaillistchecker.io/api](https://platform.emaillistchecker.io/api)
- **Email**: support@emaillistchecker.io
- **Issues**: [GitHub Issues](https://github.com/Emaillistchecker-io/emaillistchecker-go/issues)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

Made with ❤️ by [EmailListChecker](https://emaillistchecker.io)

### Batch Verification with File Upload

You can also upload CSV, TXT, or XLSX files for batch verification:

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/Emaillistchecker-io/emaillistchecker-go"
)

func main() {
	client := emaillistchecker.NewClient("your_api_key")
	
	// Upload file for batch verification
	name := "My Email List"
	batch, err := client.VerifyBatchFile("path/to/emails.csv", &name, nil, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Batch ID: %d\n", batch.ID)
	fmt.Printf("Total emails: %d\n", batch.TotalEmails)
	
	// Check progress
	for {
		status, err := client.GetBatchStatus(batch.ID)
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("Progress: %d%%\n", status.Progress)
		
		if status.Status == "completed" {
			break
		}
		
		time.Sleep(5 * time.Second)
	}
	
	// Download results
	results, err := client.GetBatchResults(batch.ID, "json", "valid")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Results: %v\n", results)
}
```

**Supported file formats:**
- CSV (.csv) - Comma-separated values
- TXT (.txt) - Plain text, one email per line
- Excel (.xlsx, .xls) - Excel spreadsheet

**File requirements:**
- Max file size: 10MB
- Max emails: 10,000 per file
- Files are automatically parsed to extract emails
