package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jram17/second-brain/services/gateway/internal/handler"
	"github.com/jram17/second-brain/services/gateway/internal/middleware"
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main(){
	//load the env
	if err:=godotenv.Load();err!=nil{
		log.Fatal("error loading the env")
		return 
	}
	authServerAdd:=os.Getenv("AUTH_SERVICE_ADDR")
	httpPort :=os.Getenv("HTTP_PORT")
	if authServerAdd== "" {
		log.Fatal("AUTH_SERVICE_ADDR is not set")
		return
	}
	if httpPort == "" {
		log.Fatal("HTTP_PORT is not set")
		return

	}
	//create a grpc connection
	conn,err:= grpc.Dial(
		authServerAdd,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err!=nil{
		log.Fatalf("failed to connect the  auth service: %v", err)
	}
	defer conn.Close()

	//create a grpc client
	authClient:=pb.NewAuthServiceClient(conn)
	authHandler:=handler.NewAuthGatewayHandler(authClient)
	//create a gin router
	r:=gin.Default()
	auth:= r.Group("/api/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login",authHandler.Login)
		auth.POST("/refresh",authHandler.RefreshToken)
	}
	
	protected:=r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authClient))
	{
		//future protected routes
	}

	//start the server
	log.Printf("Gateway service is running on port %s", httpPort)
	if err := r.Run(":" + httpPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}