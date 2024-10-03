package api

import (
	"net/http"
	"strings"
)

// isNotLogged checks if the user is logged in by verifying if the auth string is empty.
// It returns true if the user is not logged in (i.e., no bearer token is provided).
func isNotLogged(auth string) bool {

	return auth == ""
}

// validIdentifier checks if a user identifier has the correct length and format.
// It ensures that the identifier is between 3 and 16 characters long, does not contain spaces,
// and does not contain specific special characters like '?' and '_'.
func validIdentifier(identifier string) bool {
	trimmedId := strings.TrimSpace(identifier)
	return len(identifier) >= 3 && len(identifier) <= 16 && trimmedId != "" && !strings.ContainsAny(trimmedId, "?_")
}

// extractBearer extracts the bearer token from the Authorization header.
// It splits the string on spaces and returns the second part, trimming any extra whitespace.
func extractBearer(authorization string) string {
	tokens := strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}

// validateRequestingUser checks if the requesting user has a valid token for the endpoint.
// It returns 0 if valid, otherwise returns the appropriate HTTP status code.
// The function first checks if the user is logged in (i.e., a valid token is provided),
// and then ensures the token matches the user's identifier.
func validateRequestingUser(identifier string, bearerToken string) int {

	// If the requesting user has an invalid token then respond with a fobidden status
	if isNotLogged(bearerToken) {
		return http.StatusForbidden
	}

	// Validate that the identifier is in the correct format
	if !validIdentifier(identifier) {
		return http.StatusBadRequest
	}

	// Check if the requesting user's ID matches the one in the path
	if identifier != bearerToken {
		return http.StatusUnauthorized
	}
	return 0
}
