package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jram17/second-brain/services/content/internal/embedder"
	"github.com/jram17/second-brain/services/content/internal/queue"
	"github.com/jram17/second-brain/services/content/internal/vectorstore"
)

type EmbedJob struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Text      string `json:"text"`
}

func vectorID(contentID string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(contentID)).String()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	qdrantAddr := os.Getenv("QDRANT_ADDR")

	if rabbitURL == "" || qdrantAddr == "" {
		log.Fatal("missing env vars")
	}

	conn, ch, err := queue.Connect(rabbitURL)
	if err != nil {
		log.Fatalf("failed to connect to queue: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	qdrantStore, err := vectorstore.NewQdrantStore(qdrantAddr)
	if err != nil {
		log.Fatalf("failed to connect to vectorstore: %v", err)
	}
	defer qdrantStore.Close()

	msgs, err := queue.Consume(ch, "embeddings")
	if err != nil {
		log.Fatalf("failed to consume messages: %v", err)
	}
	log.Println("worker started waiting for msgs!!!")

	for msg := range msgs {
		var job EmbedJob
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("invalid msg: %v", err)
			continue
		}
		vec, err := embedder.Embed(job.Text)
		if err != nil {
			log.Printf("embedding failed: %v", err)
			continue
		}
		err = qdrantStore.Upsert(
			vectorID(job.ContentID),
			job.UserID,
			vec,
			job.Text,
		)
		if err != nil {
			log.Printf("qdrant upsert failed: %v", err)
			continue
		}

		log.Printf("processed contentId=%s", job.ContentID)
	}
}
