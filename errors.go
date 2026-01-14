package playvideo

import (
	"errors"
	"fmt"
)

// Error represents a PlayVideo API error
type Error struct {
	// HTTP status code (0 for non-HTTP errors)
	StatusCode int `json:"statusCode,omitempty"`

	// Error message
	Message string `json:"message"`

	// Error code (e.g., "authentication_error")
	Code string `json:"code,omitempty"`

	// Request ID for debugging
	RequestID string `json:"requestId,omitempty"`

	// Parameter that caused the error (for validation errors)
	Param string `json:"param,omitempty"`

	// Seconds to wait before retrying (for rate limit errors)
	RetryAfter int `json:"retryAfter,omitempty"`

	// Underlying error
	Err error `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Specific error types for type checking

// AuthenticationError represents a 401 authentication error
type AuthenticationError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *AuthenticationError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// AuthorizationError represents a 403 authorization error
type AuthorizationError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *AuthorizationError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// NotFoundError represents a 404 not found error
type NotFoundError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *NotFoundError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// ValidationError represents a 400/422 validation error
type ValidationError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *ValidationError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// ConflictError represents a 409 conflict error
type ConflictError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *ConflictError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// RateLimitError represents a 429 rate limit error
type RateLimitError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// ServerError represents a 5xx server error
type ServerError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
	Param      string
}

func (e *ServerError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("playvideo: %s (request_id: %s)", e.Message, e.RequestID)
	}
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// NetworkError represents a network/connection error
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("playvideo: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	Message string
	Err     error
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("playvideo: %s", e.Message)
}

func (e *TimeoutError) Unwrap() error {
	return e.Err
}

// WebhookSignatureError represents a webhook signature verification error
type WebhookSignatureError struct {
	Message string
}

func (e *WebhookSignatureError) Error() string {
	return fmt.Sprintf("playvideo: %s", e.Message)
}

// NewWebhookSignatureError creates a new webhook signature error
func NewWebhookSignatureError(message string) *WebhookSignatureError {
	if message == "" {
		message = "Invalid webhook signature"
	}
	return &WebhookSignatureError{
		Message: message,
	}
}

// Error type checking helpers

// IsAuthenticationError checks if the error is an authentication error
func IsAuthenticationError(err error) bool {
	var e *AuthenticationError
	return errors.As(err, &e)
}

// IsAuthorizationError checks if the error is an authorization error
func IsAuthorizationError(err error) bool {
	var e *AuthorizationError
	return errors.As(err, &e)
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var e *ValidationError
	return errors.As(err, &e)
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	var e *ConflictError
	return errors.As(err, &e)
}

// IsRateLimitError checks if the error is a rate limit error
func IsRateLimitError(err error) bool {
	var e *RateLimitError
	return errors.As(err, &e)
}

// IsServerError checks if the error is a server error
func IsServerError(err error) bool {
	var e *ServerError
	return errors.As(err, &e)
}

// IsNetworkError checks if the error is a network error
func IsNetworkError(err error) bool {
	var e *NetworkError
	return errors.As(err, &e)
}

// IsTimeoutError checks if the error is a timeout error
func IsTimeoutError(err error) bool {
	var e *TimeoutError
	return errors.As(err, &e)
}

// IsWebhookSignatureError checks if the error is a webhook signature error
func IsWebhookSignatureError(err error) bool {
	var e *WebhookSignatureError
	return errors.As(err, &e)
}

// apiErrorResponse is the structure of an API error response
type apiErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Param   string `json:"param"`
}

// parseAPIError parses an API error response and returns the appropriate error type
func parseAPIError(statusCode int, body *apiErrorResponse, requestID string) error {
	message := body.Message
	if message == "" {
		message = body.Error
	}
	if message == "" {
		message = fmt.Sprintf("HTTP %d", statusCode)
	}

	switch statusCode {
	case 401:
		return &AuthenticationError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	case 403:
		return &AuthorizationError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	case 404:
		return &NotFoundError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	case 409:
		return &ConflictError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	case 429:
		return &RateLimitError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	case 400, 422:
		return &ValidationError{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	default:
		if statusCode >= 500 {
			return &ServerError{
				StatusCode: statusCode,
				Message:    message,
				Code:       body.Code,
				RequestID:  requestID,
				Param:      body.Param,
			}
		}
		return &Error{
			StatusCode: statusCode,
			Message:    message,
			Code:       body.Code,
			RequestID:  requestID,
			Param:      body.Param,
		}
	}
}
