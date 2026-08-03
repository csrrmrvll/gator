package main

import (
	"context"
	"fmt"

	"github.com/csrrmrvll/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) < 1 {
		fmt.Println("No limit provided, defaulting to 2")
	} else {
		fmt.Sscanf(cmd.Args[0], "%d", &limit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("URL: %s\n", post.Url)
	fmt.Printf("Description: %s\n", post.Description)
	fmt.Printf("Published At: %s\n", post.PublishedAt)
	fmt.Println("-------------------------")
}
