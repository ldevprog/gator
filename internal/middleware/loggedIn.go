package middleware

import (
	"context"
	"fmt"
	"os"

	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/state"
)

func MiddlewareLoggedIn(
	handler func(s *state.State, cmd state.Command, user database.User) error,
) func(*state.State, state.Command) error {
	return func(s *state.State, cmd state.Command) error {
		currentUserName := s.Cfg.CurrentUserName
		user, err := s.DB.GetUser(context.Background(), currentUserName)
		if err != nil {
			fmt.Println("You are not logged in")
			os.Exit(1)
		}

		return handler(s, cmd, user)
	}
}
