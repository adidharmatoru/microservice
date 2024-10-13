package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type DataPaginationResponse struct {
	Data       []map[string]interface{} `json:"data"`
	Pagination PaginationResponse       `json:"pagination"`
}

// RespondWithError responds with a JSON error message
func RespondWithError(c *gin.Context, code int, message string, errors ...map[string][]string) {
	// Prepare the response map
	response := gin.H{"message": message}
	if len(errors) > 0 {
		errorMap := make(map[string][]string)
		for _, errMap := range errors {
			for field, msgs := range errMap {
				errorMap[field] = append(errorMap[field], msgs...)
			}
		}
		response["errors"] = errorMap
	}
	c.JSON(code, response)
}

// RespondWithJSON responds with a JSON payload
func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

// ValidateAndParsePagination validates pagination parameters and returns page and limit
func ValidateAndParsePagination(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	return page, limit
}

// GetPaginationLinks generates pagination links based on current page and limit
func GetPaginationLinks(c *gin.Context, page, limit, totalCount int) (nextPage, prevPage string) {
	basePath := c.Request.URL.Path
	if page*limit < totalCount {
		nextPage = fmt.Sprintf("%s?page=%d&limit=%d", basePath, page+1, limit)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("%s?page=%d&limit=%d", basePath, page-1, limit)
	}
	return nextPage, prevPage
}

// RequestOptions represents the request options for each concurrent request
type RequestOptions struct {
	Index       int                    // Index of the request
	Method      string                 // HTTP method (GET, POST, PUT, DELETE, etc.)
	URL         string                 // Request URL
	Body        interface{}            // Request body (can be JSON, string, etc.)
	QueryParams map[string]string      // Query parameters
	Headers     map[string]interface{} // Request headers
}

// GeneralResponse represents a generic response structure
type GeneralResponse struct {
	Index int         // Index of the corresponding request
	Data  interface{} // Response data
	Error string      // Error message if any
}

// PerformConcurrentRequests performs concurrent HTTP requests
func PerformConcurrentRequests(options []RequestOptions) ([]GeneralResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var responses []GeneralResponse

	for i, opt := range options {
		wg.Add(1)
		go func(i int, opt RequestOptions) {
			defer wg.Done()

			var req *http.Request
			var err error

			// Create request
			if opt.Method == "GET" {
				req, err = http.NewRequest("GET", opt.URL, nil)
			} else {
				// Serialize body to JSON if provided
				body, _ := json.Marshal(opt.Body)
				req, err = http.NewRequest(opt.Method, opt.URL, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
			}

			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				responses = append(responses, GeneralResponse{Index: i, Error: fmt.Sprintf("Failed to create request for URL %s: %v", opt.URL, err)})
				return
			}

			// Add query parameters
			q := req.URL.Query()
			for key, value := range opt.QueryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			// Add headers
			for key, value := range opt.Headers {
				req.Header.Set(key, fmt.Sprintf("%v", value))
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				mu.Lock()
				defer mu.Unlock()
				responses = append(responses, GeneralResponse{Index: i, Error: fmt.Sprintf("HTTP %s request failed for URL %s: %v", opt.Method, opt.URL, err)})
				return
			}
			defer resp.Body.Close()

			var responseData interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
				mu.Lock()
				defer mu.Unlock()
				responses = append(responses, GeneralResponse{Index: i, Error: fmt.Sprintf("Failed to decode JSON response from URL %s: %v", opt.URL, err)})
				return
			}

			mu.Lock()
			defer mu.Unlock()
			responses = append(responses, GeneralResponse{Index: i, Data: responseData})
		}(i, opt)
	}

	wg.Wait()

	// Sort responses by index
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Index < responses[j].Index
	})

	return responses, nil
}
