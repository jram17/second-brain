package handler
import (
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
	"github.com/gin-gonic/gin"
)
type AuthGatewayHandler struct{
	authClient pb.AuthServiceClient
}

type SignUpRequest struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
type RefreshTokenRequest struct{
	RefreshToken string `json:"refreshToken"`
}


//constructor
func NewAuthGatewayHandler(client pb.AuthServiceClient) *AuthGatewayHandler{
	return &AuthGatewayHandler{
		authClient: client,
	}
}


func (h *AuthGatewayHandler) Signup (c *gin.Context){
	var req SignUpRequest 
	if err := c.ShouldBindJSON(&req);err!=nil{
		c.JSON(400,gin.H{"error": err.Error()})
		return
	}
	res,err:= h.authClient.Signup(c,&pb.SignupRequest{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
	})
	if err!=nil{
		c.JSON(500,gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,gin.H{
		"accessToken":res.AccessToken,
		"refreshToken":res.RefreshToken,
	})
}

func (h *AuthGatewayHandler) Login(c *gin.Context){
	var req LoginRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(400,gin.H{"error":err.Error()})
		return
	}
	res,err:=h.authClient.Login(c,&pb.LoginRequest{
		Email: req.Email,
		Password: req.Password,
	})

	if err!=nil{
		c.JSON(500,gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,gin.H{
		"accessToken":res.AccessToken,
		"refreshToken":res.RefreshToken,
	})
}

func (h *AuthGatewayHandler) RefreshToken(c *gin.Context){
	var req RefreshTokenRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{
		c.JSON(400,gin.H{"error":err.Error()})
		return 
	}
	res,err:=h.authClient.RefreshToken(c,&pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err!=nil{
		c.JSON(500,gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,gin.H{
		"accessToken":res.AccessToken,
		"refreshToken":res.RefreshToken,
	})
}

