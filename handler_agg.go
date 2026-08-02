package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/csrrmrvll/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	fmt.Println("Collecting feeds every", time_between_reqs, "s")
	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("couldn't scrape feeds: %w", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background(), sql.NullTime{Valid: true, Time: time.Now().UTC()})
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Valid: true, Time: time.Now().UTC()},
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Item title: %s\n", item.Title)
	}
	return nil
}
