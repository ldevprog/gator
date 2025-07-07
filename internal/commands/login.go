package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("The username is required")
	}
	username := cmd.Args[0]
	user, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("This user is not registered")
		os.Exit(1)
	}
	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("The user %s has been set\n", username)

	return nil
}
