package main

import (
	"github.com/diakivm/gatorcli/internal/config"
	"github.com/diakivm/gatorcli/internal/database"
	"github.com/diakivm/gatorcli/internal/rss"
)

type state struct {
	config    *config.Config
	db        *database.Queries
	rssClient *rss.RssClient
}
