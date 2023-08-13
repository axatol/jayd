package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/axatol/jayd/pkg/config"
)

type YoutubeVideo struct {
	Kind    string `json:"kind"`
	ETag    string `json:"etag"`
	ID      string `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelID   string `json:"channelId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			Default struct {
				URL string `json:"url"`
			} `json:"default"`
		} `json:"thumbnails"`
		ChannelTitle string `json:"channelTitle"`
		CategoryID   string `json:"categoryId"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
}

type Video struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ChannelID    string `json:"channel_id"`
	ChannelTitle string `json:"channel_title"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumnail_url"`
	Duration     int64  `json:"duration"`
}

func (c *Client) Video(ctx context.Context, id string) (*Video, error) {
	if config.YoutubeAPIKey == "" {
		return nil, fmt.Errorf("youtube api key is not available")
	}

	query := url.Values{}
	query.Add("key", config.YoutubeAPIKey)
	query.Add("id", id)
	query.Add("part", "snippet,contentDetails")

	target := url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "/youtube/v3/videos",
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var parsed ListResponse[YoutubeVideo]
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Items) < 1 {
		return nil, fmt.Errorf("no results")
	}

	item := parsed.Items[0]
	duration, err := ParseDuration(item.ContentDetails.Duration)
	if err != nil {
		return nil, err
	}

	video := Video{
		ID:           item.ID,
		Title:        item.Snippet.Title,
		ChannelID:    item.Snippet.ChannelID,
		ChannelTitle: item.Snippet.ChannelTitle,
		Description:  item.Snippet.Description,
		ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
		Duration:     duration,
	}

	return &video, nil
}
