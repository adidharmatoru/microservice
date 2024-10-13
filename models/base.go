package models

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Base defines common fields and methods for all models
type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

// AdjustFieldErrors adjusts field errors to remove the model prefix and convert to lowercase
func (b *Base) AdjustFieldErrors(errMap map[string][]string, modelName string) map[string][]string {
	errors := make(map[string][]string)
	for field, msgs := range errMap {
		// Remove model prefix from field name and convert to lowercase
		key := strings.ToLower(strings.TrimPrefix(field, modelName+"."))
		errors[key] = msgs
	}
	return errors
}

// ModelName returns the name of the model
func (b *Base) ModelName() string {
	return "Base"
}

// ValidateJSONRequestAndFields validates JSON request body and struct fields
func (b *Base) ValidateJSONRequestAndFields(c *gin.Context, data interface{}) map[string][]string {
	// Parse and validate JSON request body
	if err := c.ShouldBindJSON(data); err != nil {
		// If JSON is invalid, return error detailing required fields
		errorMap := make(map[string][]string)
		// Split the error message by newline and map each line to its corresponding field
		for _, line := range strings.Split(err.Error(), "\n") {
			// Skip empty lines
			if strings.TrimSpace(line) == "" {
				continue
			}
			// Extract the field name and error message from the error line
			field := extractFieldName(line)
			msg := extractErrorMessage(line)
			if field != "" {
				// Append error message to the list of errors for the field
				errorMap[field] = append(errorMap[field], msg)
			} else {
				// If the field name cannot be extracted, add it to the "_json" key
				errorMap["_json"] = append(errorMap["_json"], line)
			}
		}
		return errorMap
	}

	// If everything is valid, return nil
	return nil
}

// extractFieldName extracts the field name from the error message
func extractFieldName(errorMessage string) string {
	// Look for substrings like "Key: 'User.Name'" or "Key: 'User.Email'"
	index := strings.Index(errorMessage, "Key: '")
	if index != -1 {
		// Extract the substring after "Key: '" and before "'"
		start := index + len("Key: '")
		end := strings.Index(errorMessage[start:], "'")
		if end != -1 {
			return errorMessage[start : start+end]
		}
	}
	return ""
}

// extractErrorMessage extracts the error message from the error string
func extractErrorMessage(errorMessage string) string {
	// Look for the "Error:" substring
	index := strings.Index(errorMessage, "Error:")
	if index != -1 {
		// Extract the message that follows "Error:"
		return strings.TrimSpace(errorMessage[index+len("Error:"):])
	}
	return ""
}
