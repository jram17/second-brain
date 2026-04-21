package handler

import (
	"context"
	"errors"

	"github.com/jram17/second-brain/services/content/internal/model"
	"github.com/jram17/second-brain/services/content/internal/scraper"
	pb "github.com/jram17/second-brain/services/content/pkg/pb"
)

type ContentHandler struct {
	pb.UnimplementedContentServiceServer
	store *model.Store
}

//constructor
func NewContentHandler (store *model.Store) *ContentHandler{
	return &ContentHandler{store: store}
}

func (h *ContentHandler) AddContent(ctx context.Context, req *pb.AddContentRequest) (*pb.AddContentResponse, error) {
	var content *model.Content
	var err error
	//content for type=url is the url
	//content for type=note is the note itself
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

func (h *ContentHandler) DeleteContent(ctx context.Context,req *pb.DeleteContentRequest)(*pb.DeleteContentResponse,error){
	err:=h.store.DeleteContentByID(ctx,req.Id,req.UserId)
	if err!=nil{
		return nil, err
	}
	return &pb.DeleteContentResponse{Success: true},nil
}