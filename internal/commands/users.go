package commands

import (
	"context"
	"fmt"

	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerUsers(s *state.State, cmd state.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		name := fmt.Sprintf("* %s", user.Name)
		if user.Name == s.Cfg.CurrentUserName {
			name = fmt.Sprintf("%s (current)", name)
		}
		fmt.Println(name)
	}

	return nil
}
