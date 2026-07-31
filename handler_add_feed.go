package main

import (
	"context"
	"fmt"

	"github.com/csrrmrvll/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Usage: addfeed <feed_url> <feed_name>")
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	// Fetch the feed to ensure it's valid
	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	name := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	// Add the feed to the database
	params := database.CreateFeedParams{
		Name:   feedName,
		Url:    feedURL,
		UserID: user.ID,
	}
	_, err = s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("couldn't add feed to database: %w", err)
	}

	fmt.Printf("Successfully added feed: %+v\n", feed)
	return nil
}
