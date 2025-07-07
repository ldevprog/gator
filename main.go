package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/levon-dalakyan/gator/internal/commands"
	"github.com/levon-dalakyan/gator/internal/config"
	"github.com/levon-dalakyan/gator/internal/database"
	"github.com/levon-dalakyan/gator/internal/middleware"
	"github.com/levon-dalakyan/gator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("An error occured reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	st := state.State{
		Cfg: &cfg,
		DB:  dbQueries,
	}
	cmds := state.Commands{
		Handlers: map[string]func(*state.State, state.Command) error{},
	}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(commands.HandlerFollow))
	cmds.Register("following", middleware.MiddlewareLoggedIn(commands.HandlerFollowing))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments were provided")
		os.Exit(1)
	}
	cmd := state.Command{
		Name: args[1],
		Args: args[2:],
	}
	err = cmds.Run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
