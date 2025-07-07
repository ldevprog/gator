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

func HandlerRegister(s *state.State, cmd state.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("The username is required")
	}
	username := cmd.Args[0]
	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("The user was created!")
	fmt.Println(user)

	return nil
}
