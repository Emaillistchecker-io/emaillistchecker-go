package emaillistchecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default API endpoint
	DefaultBaseURL = "https://platform.emaillistchecker.io/api/v1"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the EmailListChecker API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new EmailListChecker client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithConfig creates a new EmailListChecker client with custom configuration
func NewClientWithConfig(apiKey, baseURL string, timeout time.Duration) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// VerifyRequest represents a single email verification request
type VerifyRequest struct {
	Email      string `json:"email"`
	Timeout    *int   `json:"timeout,omitempty"`
	SMTPCheck  bool   `json:"smtp_check"`
}

// VerifyResponse represents a verification result
type VerifyResponse struct {
	Email        string   `json:"email"`
	Result       string   `json:"result"`
	Reason       string   `json:"reason"`
	Disposable   bool     `json:"disposable"`
	Role         bool     `json:"role"`
	Free         bool     `json:"free"`
	Score        float64  `json:"score"`
	SMTPProvider string   `json:"smtp_provider"`
	MXRecords    []string `json:"mx_records"`
	Domain       string   `json:"domain"`
	SpamTrap     bool     `json:"spam_trap"`
	MXFound      bool     `json:"mx_found"`
}

// BatchRequest represents a batch verification request
type BatchRequest struct {
	Emails      []string `json:"emails"`
	Name        string   `json:"name,omitempty"`
	CallbackURL string   `json:"callback_url,omitempty"`
	AutoStart   bool     `json:"auto_start"`
}

// BatchResponse represents a batch submission result
type BatchResponse struct {
	ID          int    `json:"id"`
	Status      string `json:"status"`
	TotalEmails int    `json:"total_emails"`
	CreatedAt   string `json:"created_at"`
}

// BatchStatusResponse represents batch status
type BatchStatusResponse struct {
	ID               int     `json:"id"`
	Status           string  `json:"status"`
	Progress         int     `json:"progress"`
	TotalEmails      int     `json:"total_emails"`
	ProcessedEmails  int     `json:"processed_emails"`
	ValidEmails      int     `json:"valid_emails"`
	InvalidEmails    int     `json:"invalid_emails"`
	UnknownEmails    int     `json:"unknown_emails"`
}

// Verify verifies a single email address
func (c *Client) Verify(email string, timeout *int, smtpCheck bool) (*VerifyResponse, error) {
	req := VerifyRequest{
		Email:     email,
		Timeout:   timeout,
		SMTPCheck: smtpCheck,
	}

	var result struct {
		Data *VerifyResponse `json:"data"`
	}

	err := c.request("POST", "/verify", req, &result)
	if err != nil {
		return nil, err
	}

	if result.Data != nil {
		return result.Data, nil
	}

	// Fallback if response doesn't have data wrapper
	var directResult VerifyResponse
	err = c.request("POST", "/verify", req, &directResult)
	return &directResult, err
}

// VerifyBatch submits emails for batch verification
func (c *Client) VerifyBatch(emails []string, name, callbackURL string, autoStart bool) (*BatchResponse, error) {
	req := BatchRequest{
		Emails:      emails,
		Name:        name,
		CallbackURL: callbackURL,
		AutoStart:   autoStart,
	}

	var result struct {
		Data *BatchResponse `json:"data"`
	}

	err := c.request("POST", "/verify/batch", req, &result)
	if err != nil {
		return nil, err
	}

	if result.Data != nil {
		return result.Data, nil
	}

	// Fallback if response doesn't have data wrapper
	var directResult BatchResponse
	err = c.request("POST", "/verify/batch", req, &directResult)
	return &directResult, err
}

// VerifyBatchFile uploads a file for batch verification (CSV, TXT, or XLSX)
func (c *Client) VerifyBatchFile(filePath string, name, callbackURL *string, autoStart bool) (*BatchResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	// Add auto_start
	if err := writer.WriteField("auto_start", strconv.FormatBool(autoStart)); err != nil {
		return nil, err
	}

	// Add optional fields
	if name != nil {
		if err := writer.WriteField("name", *name); err != nil {
			return nil, err
		}
	}
	if callbackURL != nil {
		if err := writer.WriteField("callback_url", *callbackURL); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/verify/batch/upload", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "EmailListChecker-Go/1.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		var errData map[string]interface{}
		_ = json.Unmarshal(responseBody, &errData)

		switch resp.StatusCode {
		case 401:
			msg := "Invalid API key"
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return nil, NewAuthenticationError(msg, resp.StatusCode, errData)

		case 402:
			msg := "Insufficient credits"
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return nil, NewInsufficientCreditsError(msg, resp.StatusCode, errData)

		case 422:
			msg := "Validation error"
			if errData != nil && errData["message"] != nil {
				msg = errData["message"].(string)
			}
			return nil, NewValidationError(msg, resp.StatusCode, errData)

		default:
			msg := fmt.Sprintf("API error: %d", resp.StatusCode)
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return nil, NewAPIError(msg, resp.StatusCode, errData)
		}
	}

	var result struct {
		Success bool           `json:"success"`
		Data    *BatchResponse `json:"data"`
	}

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetBatchStatus gets batch verification status
func (c *Client) GetBatchStatus(batchID int) (*BatchStatusResponse, error) {
	var result struct {
		Data *BatchStatusResponse `json:"data"`
	}

	endpoint := fmt.Sprintf("/verify/batch/%d", batchID)
	err := c.request("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	if result.Data != nil {
		return result.Data, nil
	}

	// Fallback if response doesn't have data wrapper
	var directResult BatchStatusResponse
	err = c.request("GET", endpoint, nil, &directResult)
	return &directResult, err
}

// GetBatchResults downloads batch verification results
func (c *Client) GetBatchResults(batchID int, format, filter string) (interface{}, error) {
	endpoint := fmt.Sprintf("/verify/batch/%d/results?format=%s&filter=%s", batchID, format, filter)

	var result struct {
		Data interface{} `json:"data"`
	}

	err := c.request("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// FindEmailRequest represents an email finder request
type FindEmailRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Domain    string `json:"domain"`
}

// FindEmail finds email address by name and domain
func (c *Client) FindEmail(firstName, lastName, domain string) (map[string]interface{}, error) {
	req := FindEmailRequest{
		FirstName: firstName,
		LastName:  lastName,
		Domain:    domain,
	}

	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	err := c.request("POST", "/finder/email", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// FindByDomain finds emails by domain
func (c *Client) FindByDomain(domain string, limit, offset int) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"domain": domain,
		"limit":  limit,
		"offset": offset,
	}

	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	err := c.request("POST", "/finder/domain", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// FindByCompany finds emails by company name
func (c *Client) FindByCompany(company string, limit int) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"company": company,
		"limit":   limit,
	}

	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	err := c.request("POST", "/finder/company", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetCredits gets current credit balance
func (c *Client) GetCredits() (map[string]interface{}, error) {
	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	err := c.request("GET", "/credits", nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetUsage gets API usage statistics
func (c *Client) GetUsage() (map[string]interface{}, error) {
	var result struct {
		Data map[string]interface{} `json:"data"`
	}

	err := c.request("GET", "/usage", nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetLists gets all verification lists
func (c *Client) GetLists() ([]interface{}, error) {
	var result struct {
		Data []interface{} `json:"data"`
	}

	err := c.request("GET", "/lists", nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

// DeleteList deletes a verification list
func (c *Client) DeleteList(listID int) error {
	endpoint := fmt.Sprintf("/lists/%d", listID)
	return c.request("DELETE", endpoint, nil, nil)
}

// request makes an HTTP request to the API
func (c *Client) request(method, endpoint string, body interface{}, result interface{}) error {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "EmailListChecker-Go/1.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		var errData map[string]interface{}
		_ = json.Unmarshal(respBody, &errData)

		switch resp.StatusCode {
		case 401:
			msg := "Invalid API key"
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return NewAuthenticationError(msg, resp.StatusCode, errData)

		case 402:
			msg := "Insufficient credits"
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return NewInsufficientCreditsError(msg, resp.StatusCode, errData)

		case 422:
			msg := "Validation error"
			if errData != nil && errData["message"] != nil {
				msg = errData["message"].(string)
			}
			return NewValidationError(msg, resp.StatusCode, errData)

		case 429:
			retryAfter := 60
			if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
				if val, err := strconv.Atoi(retryHeader); err == nil {
					retryAfter = val
				}
			}
			return NewRateLimitError(retryAfter, resp.StatusCode, errData)

		default:
			msg := fmt.Sprintf("API error: %d", resp.StatusCode)
			if errData != nil && errData["error"] != nil {
				msg = errData["error"].(string)
			}
			return NewAPIError(msg, resp.StatusCode, errData)
		}
	}

	// Parse successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}
