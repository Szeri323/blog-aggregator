package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/szeri323/gator/internal/config"
	"github.com/szeri323/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	registerCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registerCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registerCommands[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalf("error not enought args\n")
	}
	var s state
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error of reading config: %v\n", err)
	}
	s.cfg = &cfg

	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	cmds := commands{
		registerCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowsePosts))

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v\n", err)
	}

	// fmt.Println(config.Read())

}
