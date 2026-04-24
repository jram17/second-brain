package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jram17/second-brain/services/gateway/pkg/breaker"
	pb "github.com/jram17/second-brain/services/gateway/pkg/pb"
)

type ContentGatewayHandler struct {
	contentClient pb.ContentServiceClient
}

type AddContentRequest struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Title       string `json:"title"`
}

var contentCB = breaker.New("content-service")

func NewContentGatewayHandler(client pb.ContentServiceClient) *ContentGatewayHandler {
	return &ContentGatewayHandler{
		contentClient: client,
	}
}

func (h *ContentGatewayHandler) AddContent(c *gin.Context) {
	var req AddContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := contentCB.Execute(func() (interface{}, error) {
		return h.contentClient.AddContent(c, &pb.AddContentRequest{
			UserId:      c.GetString("userId"),
			ContentType: req.ContentType,
			Content:     req.Content,
			Title:       req.Title,
		})
	})
	if err != nil {
		c.JSON(503, gin.H{"error": "content service unavailable"})
		return
	}
	res := result.(*pb.AddContentResponse)
	c.JSON(200, gin.H{
		"content": res.Content,
	})
}

func (h *ContentGatewayHandler) GetContents(c *gin.Context) {
	result, err := contentCB.Execute(func() (interface{}, error) {
		return h.contentClient.GetContents(c, &pb.GetContentsRequest{
			UserId: c.GetString("userId"),
		})
	})
	if err != nil {
		c.JSON(503, gin.H{"error": "content service unavailable"})
		return
	}
	res := result.(*pb.GetContentsResponse)
	c.JSON(200, gin.H{
		"contents": res.Contents,
	})
}

func (h *ContentGatewayHandler) DeleteContent(c *gin.Context) {
	result, err := contentCB.Execute(func() (interface{}, error) {
		return h.contentClient.DeleteContent(c, &pb.DeleteContentRequest{
			Id:     c.Param("id"),
			UserId: c.GetString("userId"),
		})
	})
	if err != nil {
		c.JSON(503, gin.H{"error": "content service unavailable"})
		return
	}
	res := result.(*pb.DeleteContentResponse)
	c.JSON(200, gin.H{
		"success": res.Success,
	})
}
