package commands

import (
	"context"
	"fmt"

	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerUnfollow(s *state.State, cmd state.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("The url argument is required")
	}

	feedUrl := cmd.Args[0]
	feed, err := s.DB.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	err = s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
