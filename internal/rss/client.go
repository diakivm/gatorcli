package rss

import "net/http"

func NewRssClient() *RssClient {
	return &RssClient{
		client: &http.Client{},
	}
}
