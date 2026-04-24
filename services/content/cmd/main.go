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
	"github.com/jram17/second-brain/services/content/internal/queue"
	"github.com/jram17/second-brain/services/content/internal/vectorstore"
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
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//connect to mongodb
	mongouri := os.Getenv("MONGO_URI")
	qdranturi := os.Getenv("QDRANT_ADDR")
	rabbitmqurl := os.Getenv("RABBITMQ_URL")
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
	qdrantstore, err := vectorstore.NewQdrantStore(qdranturi)
	if err != nil {
		log.Fatalf("failed to connect to qdrant: %v", err)
	}
	defer qdrantstore.Close()

	// create a rabbitmq connection
	conn,ch,err:=queue.Connect(rabbitmqurl)
	if err!=nil{
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	contentHandler := handler.NewContentHandler(store,qdrantstore,ch)
	grpcServer := grpc.NewServer()
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
