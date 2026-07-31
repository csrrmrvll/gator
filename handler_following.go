package main

import (
	"context"
	"fmt"
)

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	for _, ff := range feedFollows {
		fmt.Printf("Feed: %s, User: %s\n", ff.FeedName, ff.UserName)
	}
	return nil
}
