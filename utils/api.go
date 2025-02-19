package utils

import (
	"net/http"
	"strconv"
)

// Helper function: parsePagination
func ParsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	page := ParseInt(query.Get("page"), 1)    // Default page = 1
	limit := ParseInt(query.Get("limit"), 4) // Default limit = 4
	return page, limit
}

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
