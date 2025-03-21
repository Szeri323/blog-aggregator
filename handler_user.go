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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatal("error could not call db")
	}
	printUsers(users, s)
	return nil
}

func handlerFeed(s *state, cmd command) error {
	// if len(cmd.args) == 0 {
	// 	log.Fatal("fetch command needs one argument (url)\n")
	// }
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := FetchFeed(context.Background(), feedURL)
	if err != nil {
		log.Fatal("error could not fetch the feed")
	}
	printFeed(feed)
	return nil
}
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		log.Fatal("addfeed command needs two argument (title, url)\n")
	}
	if len(cmd.args) >= 3 {
		log.Fatal("to many arguments for addfeed command it only takes two (title, url)\n")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error could not get current user: %v", err)
	}
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

func handlerFeeds(s *state, cmd command) error {
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Fatal("error follow command need one argument (url)")
	}

	feedURL := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Fatal("error could not get user from db")
	}
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

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Fatal("error could not get user from db")
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal("error could not get feed follows for user from db")
	}
	printUsersFeeds(feeds, user.Name)
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

func printFeed(feed *RSSFeed) {
	fmt.Printf(" * Title: %v\n", feed.Channel.Title)
	fmt.Printf(" * Link: %v\n", feed.Channel.Link)
	fmt.Printf(" * Description: %v\n", feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("\t- Title: %v\n", item.Title)
		fmt.Printf("\t- Description: %v\n", item.Description)
		fmt.Printf("\t- PubDate: %v\n", item.PubDate)
	}
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