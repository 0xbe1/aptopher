package aptos

import (
	"errors"
	"fmt"
)

// Common error codes returned by the Aptos API.
const (
	ErrCodeAccountNotFound  = "account_not_found"
	ErrCodeResourceNotFound = "resource_not_found"
	ErrCodeModuleNotFound   = "module_not_found"
	ErrCodeVersionPruned    = "version_pruned"
	ErrCodeInvalidInput     = "invalid_input"
	ErrCodeMempoolFull      = "mempool_is_full"
	ErrCodeVMError          = "vm_error"
	ErrCodeInternalError    = "internal_error"
)

// APIError represents an error response from the Aptos API.
type APIError struct {
	StatusCode  int     `json:"-"`
	Message     string  `json:"message"`
	ErrorCode   string  `json:"error_code"`
	VMErrorCode *uint64 `json:"vm_error_code,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("aptos api error [%s]: %s", e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("aptos api error [%d]: %s", e.StatusCode, e.Message)
}

// Is implements errors.Is for comparing API errors.
func (e *APIError) Is(target error) bool {
	var t *APIError
	if errors.As(target, &t) {
		// Match by error code if both have one
		if e.ErrorCode != "" && t.ErrorCode != "" {
			return e.ErrorCode == t.ErrorCode
		}
		// Match by status code if both have one
		if e.StatusCode != 0 && t.StatusCode != 0 {
			return e.StatusCode == t.StatusCode
		}
	}
	return false
}

// Sentinel errors for common API error conditions.
var (
	// ErrAccountNotFound is returned when the requested account does not exist.
	ErrAccountNotFound = &APIError{ErrorCode: ErrCodeAccountNotFound}

	// ErrResourceNotFound is returned when the requested resource does not exist.
	ErrResourceNotFound = &APIError{ErrorCode: ErrCodeResourceNotFound}

	// ErrModuleNotFound is returned when the requested module does not exist.
	ErrModuleNotFound = &APIError{ErrorCode: ErrCodeModuleNotFound}

	// ErrVersionPruned is returned when the requested version has been pruned.
	ErrVersionPruned = &APIError{ErrorCode: ErrCodeVersionPruned}

	// ErrInvalidInput is returned when the request input is invalid.
	ErrInvalidInput = &APIError{ErrorCode: ErrCodeInvalidInput}

	// ErrMempoolFull is returned when the mempool is full.
	ErrMempoolFull = &APIError{ErrorCode: ErrCodeMempoolFull}
)

// IsNotFound returns true if the error indicates a resource was not found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrAccountNotFound) ||
		errors.Is(err, ErrResourceNotFound) ||
		errors.Is(err, ErrModuleNotFound)
}

// IsAccountNotFound returns true if the error indicates the account was not found.
func IsAccountNotFound(err error) bool {
	return errors.Is(err, ErrAccountNotFound)
}

// IsResourceNotFound returns true if the error indicates the resource was not found.
func IsResourceNotFound(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}

// IsVersionPruned returns true if the error indicates the version was pruned.
func IsVersionPruned(err error) bool {
	return errors.Is(err, ErrVersionPruned)
}

// IsMempoolFull returns true if the error indicates the mempool is full.
func IsMempoolFull(err error) bool {
	return errors.Is(err, ErrMempoolFull)
}
