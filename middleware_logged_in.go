package main

import (
	"context"
	"log"

	"github.com/szeri323/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			log.Fatal("error could not get user from db")
		}
		return handler(s, c, user)
	}
}
