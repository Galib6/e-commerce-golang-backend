package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

func GetClaims(c *gin.Context) (*utils.JWTClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Abort()
		return nil, errors.New("missing authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.Abort()
		return nil, errors.New("invalid authorization format")
	}

	tokenString := parts[1]
	claims, err := utils.VerifyJWT(tokenString)
	if err != nil {
		c.Abort()
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// auth logic
		claims, err := GetClaims(c)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}
		// store values for handlers
		c.Set("claims", claims)
		c.Set("userId", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func IsAuthorized(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsAny, exists := c.Get("claims")
		// Roles, exists := c.Get("roles")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "unauthenticated", nil)
			c.Abort()
			return
		}
		claims := claimsAny.(*utils.JWTClaims)
		for _, role := range claims.Roles {
			if requiredRole == role {
				c.Next()
				return
			}
		}

		utils.ResponseError(c, http.StatusForbidden, "access denied", nil)
		c.Abort()
	}
}
