package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type QueryGatewayHandler struct {
	queryServiceURL string
	client          *http.Client
}

type QueryRequest struct {
	Query string `json:"query"`
}

func NewQueryGatewayHandler(queryServiceURL string) *QueryGatewayHandler {
	return &QueryGatewayHandler{
		queryServiceURL: queryServiceURL,
		client:          &http.Client{Timeout: 120 * time.Second},
	}
}

func (h *QueryGatewayHandler) Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body, err := json.Marshal(map[string]string{
		"userId": c.GetString("userId"),
		"query":  req.Query,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := h.client.Post(h.queryServiceURL+"/query", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(500, gin.H{"error": "query service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read response"})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	c.JSON(resp.StatusCode, result)
}
