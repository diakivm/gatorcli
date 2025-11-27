package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/diakivm/gatorcli/internal/database"
)

func handleLoginCommand(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gatorcli login <username>")
	}

	username := cmd.args[0]

	_, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := state.config.SetUser(username); err != nil {
		fmt.Errorf("failed to set user: %v\n", err)
	}

	log.Printf("user set: %s", username)
	return nil
}

func handleRegisterCommand(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gatorcli register <username>")
	}

	username := cmd.args[0]

	// Check if user already exists
	_, err := state.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user already exists: %s", username)
	}

	// User doesn't exist, so create them
	user, err := state.db.CreateUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	if err := state.config.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	log.Printf("user created: %+v\n", user)
	return nil
}

func handleUsersCommand(state *state, _ command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	for _, user := range users {
		if state.config.CurrentUserName == user.Name {
			log.Printf("- %s (current)\n", user.Name)
			continue
		}
		log.Printf("- %s\n", user.Name)
	}

	return nil
}

func handleResetCommand(state *state, _ command) error {
	if err := state.db.RemoveUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to remove users: %v", err)
	}

	log.Printf("users removed")
	return nil
}

func handleAggCommand(state *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := state.rssClient.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to get feed: %v", err)
	}

	data, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal feed: %v", err)
	}

	fmt.Printf("Feed: %s\n", data)

	return nil
}

func handleAddFeedCommand(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: gatorcli add-feed <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}

	usersFeed, err := state.db.CreateUsersFeed(context.Background(), database.CreateUsersFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create users feed: %v", err)
	}

	log.Printf("feed created: %+v\n", usersFeed)
	return nil
}

func handleListFeedsCommand(state *state, cmd command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}

	for _, feed := range feeds {
		log.Printf("- %s: %s (by %s)\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func handleFollowFeedCommand(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gatorcli follow <feed-url>")
	}

	feedUrl := cmd.args[0]

	feed, err := state.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to get feed by url %s: %v", feedUrl, err)
	}

	usersFeed, err := state.db.CreateUsersFeed(context.Background(), database.CreateUsersFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create users feed: %v", err)
	}

	log.Printf("followed feed %s by %s\n", usersFeed.FeedName, usersFeed.UserName)
	return nil
}

func handleListFollowedFeedsCommand(state *state, cmd command, user database.User) error {

	usersFeeds, err := state.db.GetUsersFeeds(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get users feeds: %v", err)
	}

	for _, usersFeed := range usersFeeds {
		log.Printf("- %s: %s (by %s)\n", usersFeed.FeedName, usersFeed.FeedUrl, usersFeed.UserName)
	}
	return nil
}

func handleUnfollowFeedCommand(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gatorcli unfollow <feed-url>")
	}

	feedUrl := cmd.args[0]

	feed, err := state.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to get feed by url %s: %v", feedUrl, err)
	}

	if err := state.db.RemoveUsersFeed(context.Background(), database.RemoveUsersFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("failed to remove users feed: %v", err)
	}

	log.Printf("unfollowed feed %s by %s\n", feed.Name, user.Name)
	return nil

}
