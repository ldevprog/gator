package main

import (
	"fmt"
	"log"

	"github.com/levon-dalakyan/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("An error occured reading config: %v", err)
	}

	err = cfg.SetUser("levon")
	if err != nil {
		log.Fatalf("Error setting current user: %v", err)
	}

	fmt.Println(cfg)
}
