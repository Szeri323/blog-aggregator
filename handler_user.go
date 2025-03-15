package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/szeri323/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("login command needs one argument (username)\n")
	}

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		log.Fatalf("error user does not exists\n")
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		log.Fatalf("error of setting username in state\n")
	}

	fmt.Println("Username has been set")
	return nil

}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("login command needs one argument (username)\n")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		log.Fatal("error user already exists\n")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		log.Fatal("error user cannot be created\n")
	}

	err = s.cfg.SetUser(username)

	if err != nil {
		log.Fatalf("error of setting username in state\n")
	}

	fmt.Println("Username was created in database")
	printUser(user)
	return nil

}

func handlerReset(s *state, cmd command) error {
	err := s.db.TruncateUserTable(context.Background())
	if err != nil {
		log.Fatal("error could not truncate the table")
	}
	fmt.Println("Table was truncate")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatal("error could not call db")
	}
	printUsers(users, s)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:       %v\n", user.ID)
	fmt.Printf(" * Name:     %v\n", user.Name)
}
func printUsers(users []string, s *state) {
	for _, name := range users {
		if s.cfg.CurrentUserName == name {
			fmt.Printf("* %s (current)\n", name)
		} else {
			fmt.Printf("* %s\n", name)
		}
	}
}
