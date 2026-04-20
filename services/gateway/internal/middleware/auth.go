package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthMiddleware( authClient pb.AuthServiceClient) gin.HandlerFunc{
	return func(c *gin.Context){
		//get the authorization header
		authHeader :=c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		// check the bearer format
		if !strings.HasPrefix(authHeader, "Bearer"){
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		//extract the token (stripping the Bearer prefix)
		token:=strings.TrimPrefix(authHeader, "Bearer ")
		ctx,cancel := context.WithTimeout(c.Request.Context(),5 * time.Second)
		defer cancel()

		res,err := authClient.ValidateToken(ctx,&pb.ValidateRequest{
			AccessToken: token,
		})
		if err!=nil{
			st,_:=status.FromError(err)
			switch st.Code(){
				case codes.Unauthenticated:
					c.JSON(http.StatusUnauthorized, gin.H{"error": st.Message()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service error:failed to validate token"})
			}
			c.Abort()
			return
		}

		//check if valid
		if !res.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			c.Abort()
			return
		}
		c.Set("userId", res.Userid)
		c.Next()
		
	}
}