package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/jram17/second-brain/services/content/internal/handler"
	"github.com/jram17/second-brain/services/content/internal/model"
	pb "github.com/jram17/second-brain/services/content/pkg/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//load the env
	err := godotenv.Load()
	if err!=nil{
		log.Fatal("Error loading .env file")
	}
	//connect to mongodb
	mongouri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongouri))
	if err != nil {
		log.Fatal("failed to connect to mongo :", err)
	}
	//defer the connection
	defer client.Disconnect(context.TODO())
	//ping to test
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("failed to ping mongo:", err)
	}

	//load the collection
	collection := client.Database(os.Getenv("DB_NAME")).Collection("contents")
	store := model.NewStore(collection)
	contentHandler :=handler.NewContentHandler(store)
	grpcServer :=grpc.NewServer()
	pb.RegisterContentServiceServer(grpcServer, contentHandler)
	reflection.Register(grpcServer)
	port := os.Getenv("GRPC_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen on port", port, ":", err)
	}
	fmt.Println("grpc server running on port :", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
