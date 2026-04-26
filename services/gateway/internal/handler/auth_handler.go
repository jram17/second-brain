package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"github.com/jram17/second-brain/services/gateway/pkg/breaker"
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
	"google.golang.org/grpc/status"
)

type AuthGatewayHandler struct {
	authClient pb.AuthServiceClient
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

var authCB = breaker.New("auth-service")

func NewAuthGatewayHandler(client pb.AuthServiceClient) *AuthGatewayHandler {
	return &AuthGatewayHandler{
		authClient: client,
	}
}

func handleGRPCError(c *gin.Context, err error) {
	if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
		c.JSON(503, gin.H{"error": "auth service unavailable"})
		return
	}
	if s, ok := status.FromError(err); ok {
		c.JSON(400, gin.H{"error": s.Message()})
		return
	}
	c.JSON(500, gin.H{"error": err.Error()})
}

func (h *AuthGatewayHandler) Signup(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := authCB.Execute(func() (interface{}, error) {
		return h.authClient.Signup(c, &pb.SignupRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	res := result.(*pb.SignupResponse)
	c.JSON(200, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

func (h *AuthGatewayHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := authCB.Execute(func() (interface{}, error) {
		return h.authClient.Login(c, &pb.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	res := result.(*pb.LoginResponse)
	c.JSON(200, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

func (h *AuthGatewayHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := authCB.Execute(func() (interface{}, error) {
		return h.authClient.RefreshToken(c, &pb.RefreshTokenRequest{
			RefreshToken: req.RefreshToken,
		})
	})
	if err != nil {
		handleGRPCError(c, err)
		return
	}
	res := result.(*pb.RefreshTokenResponse)
	c.JSON(200, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}
