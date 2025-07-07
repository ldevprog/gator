package commands

import (
	"context"
	"fmt"

	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerFollowing(s *state.State, cmd state.Command, user database.User) error {
	feedFollows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("You are following these feeds:")
	for _, ff := range feedFollows {
		fmt.Printf("\t- %s", ff.FeedName)
	}

	return nil
}
