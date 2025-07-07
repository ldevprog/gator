package commands

import (
	"context"
	"os"

	"github.com/levon-dalakyan/gator/internal/state"
)

func HandlerReset(s *state.State, cmd state.Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}

	return nil
}
