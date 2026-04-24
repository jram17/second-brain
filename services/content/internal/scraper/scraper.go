package scraper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/jram17/second-brain/services/content/pkg/breaker"
)

type Metadata struct {
	Title       string
	Description string
	ImageURL    string
	Content     string
}

type microlinkResponse struct {
	Status string `json:"status"`
	Data   struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Image       struct {
			URL string `json:"url"`
		} `json:"image"`
	} `json:"data"`
}

var cb = breaker.New("microlink-scraper")

func Scrape(rawURL string) (*Metadata, error) {
	if rawURL == "" {
		return nil, errors.New("empty url")
	}

	u, err := url.Parse("https://api.microlink.io")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("url", rawURL)
	u.RawQuery = q.Encode()

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	result, err := cb.Execute(func() (interface{}, error) {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result microlinkResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		if result.Status != "success" {
			return nil, errors.New("microlink: failed to scrape url")
		}

		return &Metadata{
			Title:       result.Data.Title,
			Description: result.Data.Description,
			ImageURL:    result.Data.Image.URL,
			Content:     rawURL,
		}, nil
	})
	if err != nil {
		return nil, errors.New("failed to scrape url: " + err.Error())
	}
	return result.(*Metadata), nil
}
