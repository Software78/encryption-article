package models

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // Original error, don't expose directly
} //@name APIError

type HTTPError struct {
	Code    int
	Message string
} //@name HTTPError

// Custom error handling middleware
func ErrorHandler(c *gin.Context) {
	c.Next() // Execute handlers
	fmt.Print(c.Errors)
	// Get the last error from the context
	if len(c.Errors) > 0 {
		lastError := c.Errors.Last()

		var apiErr *APIError
		if errors.As(lastError.Err, &apiErr) {
			// It's an APIError, return structured error
			c.JSON(apiErr.Code, gin.H{
				"error": apiErr, // Include the structured error
			})
			c.Abort() // Prevent further handlers from running
			return
		}

		// Handle other error types (e.g., database errors, internal errors)
		var httpErr *HTTPError
		if errors.As(lastError.Err, &httpErr) {
			c.JSON(httpErr.Code, gin.H{
				"error": httpErr.Message,
			})
			c.Abort()
			return
		}

		// Default internal server error
		log.Printf("Unhandled error: %v", lastError.Err) // Log the full error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		c.Abort()
	}
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: Code=%d, Message=%s, OriginalError=%v", e.Code, e.Message, e.Err)
}


type SuccessResponse struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
	Data    interface{}   `json:"data"`
} //@name SuccessResponse