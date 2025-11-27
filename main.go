package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/diakivm/gatorcli/internal/config"
	"github.com/diakivm/gatorcli/internal/database"
	"github.com/diakivm/gatorcli/internal/rss"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to open database: %v\n", err)
	}

	defer db.Close()

	queries := database.New(db)

	rssClient := rss.NewRssClient()

	st := &state{
		config:    &cfg,
		db:        queries,
		rssClient: rssClient,
	}

	commands := &commands{
		m: make(map[string]func(*state, command) error),
	}

	commands.register("login", handleLoginCommand)
	commands.register("register", handleRegisterCommand)
	commands.register("reset", handleResetCommand)
	commands.register("users", handleUsersCommand)
	commands.register("agg", handleAggCommand)
	commands.register("addfeed", middlewareLoggedIn(handleAddFeedCommand))
	commands.register("feeds", handleListFeedsCommand)
	commands.register("follow", middlewareLoggedIn(handleFollowFeedCommand))
	commands.register("following", middlewareLoggedIn(handleListFollowedFeedsCommand))
	commands.register("unfollow", middlewareLoggedIn(handleUnfollowFeedCommand))

	if len(os.Args) < 2 {
		log.Fatal("Usage: gatorcli <command> <args...>")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := commands.run(st, cmd); err != nil {
		log.Fatalf("failed to run command: %v\n", err)
	}

}
