package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerAddFeed(s *state.State, cmd state.Command) error {
	if len(cmd.Args) < 2 {
		fmt.Println("Not enough arguments were provided. Expected to get args for feed name and feed url")
		os.Exit(1)
	}

	currentUserName := s.Cfg.CurrentUserName
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	user, err := s.DB.GetUser(context.Background(), currentUserName)
	if err != nil {
		fmt.Println("You are not logged in")
		os.Exit(1)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Println("Error creating feed")
		os.Exit(1)
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		fmt.Println("Error creating feed_follow")
		os.Exit(1)
	}

	fmt.Println(feed)

	return nil
}
