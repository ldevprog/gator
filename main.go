package main

import (
	"fmt"
	"log"
	"os"

	"github.com/levon-dalakyan/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, found := c.handlers[cmd.name]
	if !found {
		return fmt.Errorf("Command not found")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("The username is required")
	}
	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("The user %s has been set\n", username)

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("An error occured reading config: %v", err)
	}

	st := state{
		config: &cfg,
	}
	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments were provided")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
