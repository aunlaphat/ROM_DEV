package api

import (
	"net/http"
	"strconv"
)

// Helper function: parsePagination
func parsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	page := parseInt(query.Get("page"), 1)    // Default page = 1
	limit := parseInt(query.Get("limit"), 10) // Default limit = 10
	return page, limit
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

