package coolify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ConflictDetail represents detailed resource conflict information returned on 409 Conflict responses.
type ConflictDetail struct {
	Domain       string  `json:"domain"`
	ResourceName string  `json:"resource_name"`
	ResourceUUID *string `json:"resource_uuid,omitempty"`
	ResourceType string  `json:"resource_type"`
	Message      string  `json:"message"`
}

// Error represents an error returned by the Coolify API.
// It parses structured Laravel-style validation errors (422) and Traefik/Caddy domain conflicts (409).
type Error struct {
	StatusCode int                         `json:"-"`
	Message    string                      `json:"message"`
	Warning    string                      `json:"warning,omitempty"`
	Errors     map[string][]string         `json:"errors,omitempty"`
	Conflicts  []ConflictDetail            `json:"conflicts,omitempty"`
	RawBody    string                      `json:"-"`
}

// Error formats the structured Coolify error into a highly readable error string.
func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("coolify: api error (status %d)", e.StatusCode))

	if e.Message != "" {
		sb.WriteString(": ")
		sb.WriteString(e.Message)
	}

	if e.Warning != "" {
		sb.WriteString(" (warning: ")
		sb.WriteString(e.Warning)
		sb.WriteString(")")
	}

	// Format Laravel validation errors
	if len(e.Errors) > 0 {
		var validationErrs []string
		for field, messages := range e.Errors {
			validationErrs = append(validationErrs, fmt.Sprintf("%s: %s", field, strings.Join(messages, ", ")))
		}
		sb.WriteString(" [validation errors: ")
		sb.WriteString(strings.Join(validationErrs, "; "))
		sb.WriteString("]")
	}

	// Format domain conflicts
	if len(e.Conflicts) > 0 {
		var conflictDetails []string
		for _, conflict := range e.Conflicts {
			conflictDetails = append(conflictDetails, fmt.Sprintf("%s conflict with %s (%s)", conflict.Domain, conflict.ResourceName, conflict.ResourceType))
		}
		sb.WriteString(" [conflicts: ")
		sb.WriteString(strings.Join(conflictDetails, "; "))
		sb.WriteString("]")
	}

	return sb.String()
}

// checkError checks an HTTP response and converts any non-2xx status code into a parsed Error.
func parseResponseError(statusCode int, body []byte) error {
	apiErr := &Error{
		StatusCode: statusCode,
		RawBody:    string(body),
	}

	// Try to unmarshal JSON response
	if len(body) > 0 && json.Unmarshal(body, apiErr) == nil {
		// Ensure a default message is present if parsed JSON is empty
		if apiErr.Message == "" {
			apiErr.Message = fmt.Sprintf("HTTP status %d", statusCode)
		}
		return apiErr
	}

	// Fallback to plain text if JSON unmarshaling fails
	rawMessage := strings.TrimSpace(string(body))
	if strings.Contains(strings.ToLower(rawMessage), "<html") {
		// If it's an HTML page, do not dump the raw HTML. Just say Unauthenticated or HTML response.
		apiErr.Message = fmt.Sprintf("HTML response (HTTP status %d)", statusCode)
		if statusCode == http.StatusUnauthorized {
			apiErr.Message = "Unauthenticated (401)"
		} else if statusCode == http.StatusForbidden {
			apiErr.Message = "Forbidden (403)"
		} else if statusCode == http.StatusNotFound {
			apiErr.Message = "Not Found (404)"
		}
	} else {
		apiErr.Message = rawMessage
	}

	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("HTTP status %d", statusCode)
	}
	return apiErr
}
