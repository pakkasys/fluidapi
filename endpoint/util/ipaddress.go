package util

import (
	"net/http"
	"strings"
)

func RequestIPAddress(request *http.Request) string {
	// Check for the IP in X-Forwarded-For header
	forwarded := request.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IP addresses; take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to using RemoteAddr if X-Forwarded-For is not available
	ip := request.RemoteAddr
	// RemoteAddr contains IP:Port, so split by ':' and take the first part
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		return ip[:colonIndex]
	}
	return ip
}
