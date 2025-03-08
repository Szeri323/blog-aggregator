package main

import (
	"fmt"
	"log"

	"github.com/szeri323/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("login command needs one argument (username)\n")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		log.Fatalf("error of setting username in state\n")
	}

	fmt.Println("Username has been set")
	return nil

}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error of reading config: %v\n", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("szeri")
	if err != nil {
		log.Fatal("error of setting user config: %v\n", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal("error of reading config: %v\n", err)
	}

	var s state
	s.cfg = &cfg
	var cmd command
	cmd.name = "login"
	cmd.args = []string{"test"}
	handlerLogin(&s, cmd)

	fmt.Println(cfg)
}
