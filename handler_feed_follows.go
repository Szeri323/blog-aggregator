package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/szeri323/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		log.Fatal("error follow command need one argument (url)")
	}

	feedURL := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		log.Fatal("error could not get the feed from db")
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal("error could not create feed follow")
	}

	printFeedFollow(user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal("error could not get feed follows for user from db")
	}
	printUsersFeeds(feeds, user.Name)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		log.Fatal("error unfollow command need one argument (url)")
	}
	feedURL := cmd.args[0]
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Url:    feedURL,
		UserID: user.ID,
	})
	if err != nil {
		log.Fatal("error could not unfollow the feed")
	}
	return nil
}

func printFeedFollow(user_name string, feed_name string) {
	fmt.Println("Feed follow created.")
	fmt.Println(user_name)
	fmt.Println(feed_name)
}

func printUsersFeeds(feeds []database.GetFeedFollowsForUserRow, user_name string) {
	fmt.Println(user_name)
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
}
