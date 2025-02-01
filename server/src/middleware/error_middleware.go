package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Custom error types
type (
	ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	AppError struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	}
)

// Implement error interface for ValidationError
func (e *ValidationError) Error() string {
	return e.Message
}

// Implement error interface for AppError
func (e *AppError) Error() string {
	return e.Message
}

// ErrorHandler middleware for comprehensive error management
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			handleErrors(c)
			return
		}
	}
}

func handleErrors(c *gin.Context) {
	// Collect all errors
	var validationErrors []ValidationError
	var appError *AppError
	var primaryError error

	for _, e := range c.Errors {
		switch err := e.Err.(type) {
		case *AppError:
			appError = err
		case *ValidationError:
			validationErrors = append(validationErrors, *err)
		default:
			primaryError = e.Err
		}
	}

	// Prioritize error handling
	switch {
	case appError != nil:
		// Custom application error
		c.JSON(appError.Code, appError)
		return

	case len(validationErrors) > 0:
		// Validation errors
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed",
			"details": validationErrors,
		})
		return

	case primaryError != nil:
		// Handle specific known error types
		switch {
		case errors.Is(primaryError, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Resource Not Found",
				
			})
		case errors.Is(primaryError, gorm.ErrDuplicatedKey):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Duplicate Record",
			})
		case errors.Is(primaryError, bcrypt.ErrMismatchedHashAndPassword):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			case errors.Is(primaryError,  errors.New("field is not encryped")):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "field is not encryped",
			})
		default:
			// Log unexpected errors
			log.Printf("Unhandled error: %v", primaryError)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}
}

// Helper functions for creating specific error types
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewAppError(code int, message string, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
