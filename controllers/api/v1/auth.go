package v1

import (
	"net/http"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

// GenerateJWT godoc
// @Summary Generate a new JWT token
// @Description Generates a new JWT token for user authentication and authorization
// @Tags authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} object "token: <generated_token>"
// @Failure 500 {object} object "message: Failed to generate JWT token"
// @Router /oauth/token [post]
func GenerateJWT(c *gin.Context) {
	// Define your JWT claims
	claims := jwt.MapClaims{
		"sub":   "dummyuser",
		"iss":   os.Getenv("AUTH0_DOMAIN"),
		"aud":   os.Getenv("AUTH0_AUDIENCE"),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),               // Token expires after 24 hours
		"scope": "create:users read:users update:users delete:users", // Define the scope of the token for CRUD operations on users
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret
	secret := os.Getenv("AUTH0_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate JWT token"})
		return
	}

	// Return the generated JWT token
	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
