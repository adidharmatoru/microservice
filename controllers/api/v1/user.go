package v1

import (
	"microservice/controllers/api"
	"microservice/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Define UsersResponse struct
type UsersResponse struct {
	Data       []models.User          `json:"data"`
	Pagination api.PaginationResponse `json:"pagination"`
}

// OptionsUsers handles OPTIONS requests for the /users endpoint
func OptionsUsers(c *gin.Context) {
	c.Header("Allow", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
	c.Header("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	api.RespondWithJSON(c, http.StatusOK, gin.H{})
}

// HeadUsers handles HEAD requests for the /users endpoint
func HeadUsers(c *gin.Context) {
	// Add any necessary logic here
	api.RespondWithJSON(c, http.StatusOK, gin.H{})
}

// ListUsers godoc
// @Summary Get all users
// @Description Get all users with optional filtering and pagination
// @Tags users
// @Accept  json
// @Produce  json
// @Param name query string false "Name"
// @Param age query int false "Age"
// @Param ids query string false "Comma-separated list of user IDs"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} UsersResponse
// @Router /users [get]
func ListUsers(c *gin.Context) {
	var users []models.User
	page, limit := api.ValidateAndParsePagination(c)
	query := models.DB

	// Optional filtering
	if name := c.Query("name"); name != "" {
		query = query.Where("name = ?", name)
	}
	if age, err := strconv.Atoi(c.Query("age")); err == nil {
		query = query.Where("age = ?", age)
	}
	if ids := c.Query("ids"); ids != "" {
		idList := strings.Split(ids, ",")
		query = query.Where("id IN (?)", idList)
	}

	// Apply limit and offset
	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	// Execute query
	query.Find(&users)

	// Calculate total count
	var totalCount int
	query.Model(&models.User{}).Count(&totalCount)

	// Create pagination metadata
	nextPage, prevPage := api.GetPaginationLinks(c, page, limit, totalCount)

	// Create response object
	response := UsersResponse{
		Data:       users,
		Pagination: api.PaginationResponse{Next: nextPage, Previous: prevPage, Total: totalCount},
	}

	api.RespondWithJSON(c, http.StatusOK, response)
}

// GetUser godoc
// @Summary Get a single user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")

	// Convert user ID to integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Retrieve user from the database
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		api.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Return the user as JSON response
	api.RespondWithJSON(c, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User"
// @Security BearerToken
// @Success 200 {object} models.User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	// Validate JSON request body
	if errMap := user.ValidateJSONRequestAndFields(c, &user); len(errMap) > 0 {
		errors := user.AdjustFieldErrors(errMap)
		api.RespondWithError(c, http.StatusBadRequest, "Invalid JSON format", errors)
		return
	}

	// Create user in the database
	if err := models.DB.Create(&user).Error; err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Return the created user as JSON response
	api.RespondWithJSON(c, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user by id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Security BearerToken
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Check if user exists
	if err := models.DB.First(&user, id).Error; err != nil {
		api.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Validate JSON request body
	if err := user.ValidateJSONRequestAndFields(c, &user); err != nil {
		api.RespondWithError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Save updated user to the database
	if err := models.DB.Save(&user).Error; err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Return the updated user as JSON response
	api.RespondWithJSON(c, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Security BearerToken
// @Success 204
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.User{}, id).Error; err != nil {
		api.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}
	c.Status(http.StatusNoContent)
}

// DummyListUsers godoc
// @Summary Test goroutine to fetch users
// @Description Fetches users concurrently from dummy API with pagination
// @Tags users
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} api.PaginationResponse
// @Router /users/dummy [get]
func DummyListUsers(c *gin.Context) {
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Prepare options for concurrent requests
	var options []api.RequestOptions
	for i := 0; i < limit; i++ {
		options = append(options, api.RequestOptions{
			Method: "GET",
			URL:    "https://dummyjson.com/users/" + strconv.Itoa(i+1),
		})
	}

	// Perform concurrent requests
	responses, err := api.PerformConcurrentRequests(options)
	if err != nil {
		api.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Extract user data from responses
	var users []map[string]interface{}
	for _, response := range responses {
		if response.Error != "" {
			users = append(users, map[string]interface{}{"error": response.Error})
		} else {
			userData := response.Data.(map[string]interface{})
			users = append(users, userData)
		}
	}

	// Dummy total count
	totalCount := 100

	// Calculate pagination metadata
	nextPage, prevPage := api.GetPaginationLinks(c, page, limit, totalCount)

	response := api.DataPaginationResponse{
		Data: users,
		Pagination: api.PaginationResponse{
			Next:     nextPage,
			Previous: prevPage,
			Total:    totalCount,
		},
	}
	api.RespondWithJSON(c, http.StatusOK, response)
}
