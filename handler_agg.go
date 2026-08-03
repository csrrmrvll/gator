package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/csrrmrvll/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: parsePubDate(item.PubDate),
		})
		if err != nil {
			log.Printf("Couldn't create post %s: %v", item.Title, err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func parsePubDate(pubDate string) sql.NullTime {
	t, err := time.Parse(time.RFC1123Z, pubDate)
	if err != nil {
		t, err = time.Parse(time.RFC1123, pubDate)
		if err != nil {
			log.Printf("Couldn't parse pubDate %s: %v", pubDate, err)
			return sql.NullTime{Valid: false}
		}
	}
	return sql.NullTime{Time: t, Valid: true}
}
