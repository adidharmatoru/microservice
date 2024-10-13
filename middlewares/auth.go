package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtMiddleware *jwtmiddleware.JWTMiddleware

func init() {
	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			aud := os.Getenv("AUTH0_AUDIENCE")
			if !token.Claims.(jwt.MapClaims).VerifyAudience(aud, false) {
				return nil, jwt.NewValidationError("invalid audience", jwt.ValidationErrorAudience)
			}
			iss := os.Getenv("AUTH0_DOMAIN")
			if !token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false) {
				return nil, jwt.NewValidationError("invalid issuer", jwt.ValidationErrorIssuer)
			}
			if !token.Claims.(jwt.MapClaims).VerifyExpiresAt(time.Now().Unix(), false) {
				return nil, jwt.NewValidationError("token expired", jwt.ValidationErrorExpired)
			}

			secret := []byte(os.Getenv("AUTH0_SECRET"))
			return secret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := jwtMiddleware.CheckJWT(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "details": err.Error()})
			c.Abort()
			return
		}
		// Extract the token from the request context
		token, _ := c.Request.Context().Value("user").(*jwt.Token)
		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Set the token in the Gin context
		c.Set("user", token)
		c.Next()
	}
}

// CheckScope is a middleware to check if the JWT token has the required scope.
func CheckScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the context
		token, exists := c.Get("user")
		log.Println("token: ", token)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Assert the token to the correct type
		jwtToken, ok := token.(*jwt.Token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract claims and check for scope
		claims := jwtToken.Claims.(jwt.MapClaims)
		scopes, exists := claims["scope"].(string)
		if !exists || !strings.Contains(scopes, scope) {
			c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient scope"})
			c.Abort()
			return
		}

		c.Next()
	}
}
