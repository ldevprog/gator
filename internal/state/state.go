package state

import (
	"fmt"

	"github.com/levon-dalakyan/gator/internal/config"
	"github.com/levon-dalakyan/gator/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, found := c.Handlers[cmd.Name]
	if !found {
		return fmt.Errorf("Command not found")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}
