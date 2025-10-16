package emaillistchecker

import "fmt"

// Error represents a base error from the EmailListChecker API
type Error struct {
	Message      string
	StatusCode   int
	ResponseData map[string]interface{}
}

func (e *Error) Error() string {
	return e.Message
}

// AuthenticationError is returned when API authentication fails
type AuthenticationError struct {
	*Error
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string, statusCode int, responseData map[string]interface{}) *AuthenticationError {
	return &AuthenticationError{
		Error: &Error{
			Message:      message,
			StatusCode:   statusCode,
			ResponseData: responseData,
		},
	}
}

// InsufficientCreditsError is returned when account has insufficient credits
type InsufficientCreditsError struct {
	*Error
}

// NewInsufficientCreditsError creates a new insufficient credits error
func NewInsufficientCreditsError(message string, statusCode int, responseData map[string]interface{}) *InsufficientCreditsError {
	return &InsufficientCreditsError{
		Error: &Error{
			Message:      message,
			StatusCode:   statusCode,
			ResponseData: responseData,
		},
	}
}

// RateLimitError is returned when API rate limit is exceeded
type RateLimitError struct {
	*Error
	RetryAfter int
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(retryAfter int, statusCode int, responseData map[string]interface{}) *RateLimitError {
	return &RateLimitError{
		Error: &Error{
			Message:      fmt.Sprintf("Rate limit exceeded. Retry after %d seconds", retryAfter),
			StatusCode:   statusCode,
			ResponseData: responseData,
		},
		RetryAfter: retryAfter,
	}
}

// ValidationError is returned when request validation fails
type ValidationError struct {
	*Error
}

// NewValidationError creates a new validation error
func NewValidationError(message string, statusCode int, responseData map[string]interface{}) *ValidationError {
	return &ValidationError{
		Error: &Error{
			Message:      message,
			StatusCode:   statusCode,
			ResponseData: responseData,
		},
	}
}

// APIError is returned for general API errors
type APIError struct {
	*Error
}

// NewAPIError creates a new API error
func NewAPIError(message string, statusCode int, responseData map[string]interface{}) *APIError {
	return &APIError{
		Error: &Error{
			Message:      message,
			StatusCode:   statusCode,
			ResponseData: responseData,
		},
	}
}
