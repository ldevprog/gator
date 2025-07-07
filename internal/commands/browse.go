package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerBrowse(s *state.State, cmd state.Command, user database.User) error {
	limit := 2
	var err error
	if len(cmd.Args) > 0 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("The provided limit argument is not integer")
		}
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("- %s\n", p.Title)
	}

	return nil
}
