package handler

import (
	"net/http"
	"strings"

	"1litw/application"
	"1litw/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware enforces that a user must be authenticated.
// If the token is missing or invalid, it aborts the request.
func AuthMiddleware(jwtSecret string, userUseCase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := extractUserFromToken(c, jwtSecret, userUseCase)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuthMiddleware tries to authenticate a user, but does not fail if authentication fails.
// If a token is present and valid, the user is set in the context.
// If no token is present or the token is invalid, it fetches the 'anonymous' user and sets it in the context.
func OptionalAuthMiddleware(jwtSecret string, userUseCase *application.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := extractUserFromToken(c, jwtSecret, userUseCase)
		if err != nil {
			// Token is missing or invalid, proceed with the anonymous user.
			anonymousUser, err := userUseCase.GetAnonymousUser(c.Request.Context())
			if err != nil {
				// This is a server error, as the anonymous user should always exist.
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not load anonymous user"})
				return
			}
			c.Set("user", anonymousUser)
			c.Next()
			return
		}

		// Successfully authenticated user
		c.Set("user", user)
		c.Next()
	}
}

// extractUserFromToken contains the logic to get a user from a JWT in the request.
// It checks cookies and the Authorization header.
func extractUserFromToken(c *gin.Context, jwtSecret string, userUseCase *application.UserUseCase) (*domain.User, error) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, &AuthError{"authorization token not provided"}
		}
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if tokenString == "" {
		return nil, &AuthError{"authorization token not provided"}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, &AuthError{"invalid token"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &AuthError{"invalid token claims"}
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, &AuthError{"invalid user id in token"}
	}
	userID := int64(userIDFloat)

	user, err := userUseCase.GetMe(c.Request.Context(), userID)
	if err != nil || user == nil {
		return nil, &AuthError{"user not found"}
	}

	return user, nil
}

// AuthError is a custom error type for authentication failures.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

