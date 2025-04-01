package card

import "fmt"

// CardError represents an error that occurred during card operations
type CardError struct {
	Code    CardErrorCode
	Message string
}

// CardErrorCode represents the type of card error
type CardErrorCode int

const (
	// ErrCardNotFound indicates that a card template was not found
	ErrCardNotFound CardErrorCode = iota
)

// Error implements the error interface
func (e *CardError) Error() string {
	return fmt.Sprintf("card error: %s", e.Message)
}

// NewCardError creates a new card error
func NewCardError(code CardErrorCode, message string) *CardError {
	return &CardError{
		Code:    code,
		Message: message,
	}
} 