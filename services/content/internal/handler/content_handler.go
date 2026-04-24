package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jram17/second-brain/services/content/internal/model"
	"github.com/jram17/second-brain/services/content/internal/queue"
	"github.com/jram17/second-brain/services/content/internal/scraper"
	"github.com/jram17/second-brain/services/content/internal/vectorstore"
	pb "github.com/jram17/second-brain/services/content/pkg/pb"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ContentHandler struct {
	pb.UnimplementedContentServiceServer
	store       *model.Store
	vectorStore *vectorstore.QdrantStore
	mqChannel   *amqp.Channel
}

type EmbedMessage struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Text      string `json:"text"`
}

func NewContentHandler(store *model.Store, vs *vectorstore.QdrantStore, mqChannel *amqp.Channel) *ContentHandler {
	return &ContentHandler{store: store, vectorStore: vs, mqChannel: mqChannel}
}

func vectorID(contentID string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(contentID)).String()
}

func (h *ContentHandler) AddContent(ctx context.Context, req *pb.AddContentRequest) (*pb.AddContentResponse, error) {
	var content *model.Content
	var err error

	switch req.ContentType {
	case "note":
		content, err = h.store.AddContent(ctx, &model.Content{
			UserID:      req.UserId,
			ContentType: req.ContentType,
			Title:       req.Title,
			Content:     req.Content,
		})
	case "link":
		meta, err := scraper.Scrape(req.Content)
		if err != nil {
			return nil, err
		}
		content, err = h.store.AddContent(ctx, &model.Content{
			UserID:      req.UserId,
			ContentType: req.ContentType,
			Title:       meta.Title,
			Description: meta.Description,
			ImageURL:    meta.ImageURL,
			Content:     meta.Content,
		})
	default:
		return nil, errors.New("invalid content type")
	}
	if err != nil {
		return nil, err
	}

	// publish to queue for async embedding
	text := content.Title + " " + content.Description + " " + content.Content

	msg, err := json.Marshal(EmbedMessage{
		ContentID: content.ID.Hex(),
		UserID:    content.UserID,
		Text:      text,
	})
	if err != nil {
		log.Printf("err during marshalling data: %v", err)
	}
	if err == nil {
		if pubErr := queue.Publish(h.mqChannel, "embeddings", msg); pubErr != nil {
			log.Printf("error publishing to queue: %v", pubErr)
		}
	}
	return &pb.AddContentResponse{
		Content: &pb.Content{
			Id:          content.ID.Hex(),
			UserId:      content.UserID,
			ContentType: content.ContentType,
			Title:       content.Title,
			Description: content.Description,
			ImageUrl:    content.ImageURL,
			Content:     content.Content,
			CreatedAt:   content.CreatedAt,
		},
	}, nil
}

func (h *ContentHandler) GetContents(ctx context.Context, req *pb.GetContentsRequest) (*pb.GetContentsResponse, error) {
	contents, err := h.store.GetContentByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	var pbContents []*pb.Content
	for _, c := range contents {
		pbContents = append(pbContents, &pb.Content{
			Id:          c.ID.Hex(),
			UserId:      c.UserID,
			ContentType: c.ContentType,
			Title:       c.Title,
			Description: c.Description,
			ImageUrl:    c.ImageURL,
			Content:     c.Content,
			CreatedAt:   c.CreatedAt,
		})
	}
	return &pb.GetContentsResponse{Contents: pbContents}, nil
}

func (h *ContentHandler) DeleteContent(ctx context.Context, req *pb.DeleteContentRequest) (*pb.DeleteContentResponse, error) {
	err := h.store.DeleteContentByID(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	// delete from qdrant directly (sync is fine for delete)
	h.vectorStore.Delete(vectorID(req.Id))
	return &pb.DeleteContentResponse{Success: true}, nil
}
