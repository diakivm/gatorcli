package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

func (c *RssClient) GetFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	result := &RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "gatorcli/1.0")

	res, err := c.client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	if err := xml.NewDecoder(res.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}

	result.Channel.Title = html.UnescapeString(result.Channel.Title)
	result.Channel.Description = html.UnescapeString(result.Channel.Description)
	for index, item := range result.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		result.Channel.Item[index] = item
	}

	return result, nil
}
