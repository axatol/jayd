package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/axatol/jayd/config"
)

type Video struct {
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

func (c *Client) Video(ctx context.Context, id string) (*Video, error) {
	if config.YoutubeAPIKey == "" {
		return nil, fmt.Errorf("youtube api key is not available")
	}

	target := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?key=%s&id=%s&part=%s",
		config.YoutubeAPIKey,
		id,
		"snippet,contentDetails",
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, http.NoBody)
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

	var parsed ListResponse[Video]
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Items) < 1 {
		return nil, fmt.Errorf("no results")
	}

	return &parsed.Items[0], nil
}
