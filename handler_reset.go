package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.TruncateUsersTable(context.Background())
	if err != nil {
		log.Fatal("error could not truncate the users table")
	}
	err = s.db.TruncateFeedsTable(context.Background())
	if err != nil {
		log.Fatal("error could not truncate the feeds table")
	}
	err = s.db.TruncateFeedFollowsTable(context.Background())
	if err != nil {
		log.Fatal("error could not truncate the feed follows table")
	}
	fmt.Println("Tables were truncate")
	return nil
}
