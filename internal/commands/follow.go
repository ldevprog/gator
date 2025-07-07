package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerFollow(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("The url argument is required")
	}

	feedUrl := cmd.Args[0]
	currentUser, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feedToFollow, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feedToFollow.ID,
	})
	if err != nil {
		return err
	}

	out := fmt.Sprintf("User %s is now following %s feed", feedFollow.UserName, feedFollow.FeedName)
	fmt.Println(out)

	return nil
}
