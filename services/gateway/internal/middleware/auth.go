package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jram17/second-brain/services/gateway/pkg/breaker"
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
)

var authMiddlewareCB = breaker.New("auth-validate")

func AuthMiddleware(authClient pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		result, err := authMiddlewareCB.Execute(func() (interface{}, error) {
			return authClient.ValidateToken(ctx, &pb.ValidateRequest{
				AccessToken: token,
			})
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
			c.Abort()
			return
		}

		res := result.(*pb.ValidateResponse)
		if !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			c.Abort()
			return
		}
		c.Set("userId", res.Userid)
		c.Next()
	}
}
