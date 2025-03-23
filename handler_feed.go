package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/szeri323/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		log.Fatal("addfeed command needs two argument (title, url)\n")
	}
	if len(cmd.args) >= 3 {
		log.Fatal("to many arguments for addfeed command it only takes two (title, url)\n")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeeds(context.Background(), database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal("error could not create feed follow")
	}

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		log.Fatal("feeds not need any additional arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatal("error could not get feeds")
	}
	printFeeds(feeds)
	return nil
}

func printFeeds(feeds []database.GetFeedsRow) {
	for _, feed := range feeds {
		fmt.Println("*********")
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2)
		fmt.Println("*********")
	}
}
