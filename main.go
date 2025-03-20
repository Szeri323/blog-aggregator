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
	commandMap map[string]func(*state, command) error
}

func (c *commands) login(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
func (c *commands) reset(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
func (c *commands) users(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func (c *commands) agg(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
func (c *commands) addfeed(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}
func (c *commands) feeds(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.commandMap[cmd.name]
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

	var cmds commands
	cmds.commandMap = make(map[string]func(*state, command) error)
	cmds.login("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.reset("reset", handlerReset)
	cmds.users("users", handlerUsers)
	cmds.agg("agg", handlerFeed)
	cmds.addfeed("addfeed", handlerAddFeed)
	cmds.feeds("feeds", handlerFeeds)

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v\n", err)
	}

	// fmt.Println(config.Read())

}
