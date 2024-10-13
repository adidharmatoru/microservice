package main

import (
	"log"
	"microservice/database"
	_ "microservice/docs"
	"microservice/routes"

	"github.com/joho/godotenv"
)

// @title Passport Auth API
// @version 1.0
// @description Microservice for Passport Auth API

// @host passport.adidharmatoru.dev
// @BasePath /api/v1

// @securityDefinitions.BearerToken
// @type apiKey
// @name Authorization
// @in header
// @description Enter the JWT token with "Bearer " prefix

// @securityDefinitions.oauth2
// @type oauth2
// @tokenUrl /api/v1/oauth/token
// @flow password
// @scopes.read Access to read endpoints
// @scopes.write Access to write endpoints
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	database.ConnectDatabase()

	// Setup the router
	router := routes.SetupRouter()

	// Run the server
	router.Run(":8080")
}
