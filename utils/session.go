package utils

import (
	"net"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
)

// generateSessionToken creates a new session token using a UUID version 4.
//
// This function generates a unique session token by creating a new UUID version 4.
// It's typically used for creating session identifiers in web applications.
//
// Returns:
//   - string: A string representation of the generated UUID, to be used as a session token.
//   - error: An error if the UUID generation fails, or nil if successful.
func GenerateSessionToken() (string, error) {
	u, err := uuid.NewV4() // Generate a UUID version 4 for the session token
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// getIP extracts the client's IP address from an HTTP request.
//
// This function attempts to determine the client's IP address by checking various
// HTTP headers and the request's remote address. It prioritizes the following
// sources in order:
//  1. X-Forwarded-For header
//  2. X-Real-IP header
//  3. RemoteAddr from the request
//
// Parameters:
//   - r: An *http.Request object representing the incoming HTTP request.
//
// Returns:
//   - string: The extracted IP address as a string. If no valid IP address is
//     found in the headers, it returns the remote address from the request.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
