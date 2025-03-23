package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/szeri323/gator/internal/database"
)

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("error could not get next feed: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		log.Fatalf("error could not update the feed: %v", err)
	}

	RSSFeed, err := FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Fatalf("error could not fetch the RSSFeed: %v", err)
	}
	// printRSSFeed(RSSFeed)

	/* Use RSSFeed Item isted of RSSFeed */
	for _, RSSItem := range RSSFeed.Channel.Item {
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       RSSItem.Title,
			Url:         RSSItem.Link,
			Description: RSSItem.Description,
			PublishedAt: RSSItem.PubDate,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates") {
				fmt.Println("DUPLICATE URL detected, skipping...")
				continue
			}
			log.Fatalf("error could not create post in db: %v", err)
		}
		fmt.Println("post created:")
		fmt.Println(post)
	}

}

func handlerAgg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration("1s")
	if err != nil {
		log.Fatal("error could not parse time")
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

// func printRSSFeed(feed *RSSFeed) {
// 	fmt.Printf(" * Title: %v\n", feed.Channel.Title)
// 	fmt.Printf(" * Link: %v\n", feed.Channel.Link)
// 	fmt.Printf(" * Description: %v\n", feed.Channel.Description)
// 	for _, item := range feed.Channel.Item {
// 		fmt.Printf("\t- Title: %v\n", item.Title)
// 		fmt.Printf("\t- Description: %v\n", item.Description)
// 		fmt.Printf("\t- PubDate: %v\n", item.PubDate)
// 	}
// }

// func printFeed(feed database.Feed) {
// 	fmt.Printf("%s\n", feed.ID)
// 	fmt.Printf("%s\n", feed.Name)
// 	fmt.Printf("%s\n", feed.Url)
// 	fmt.Printf("%s\n", feed.CreatedAt)
// 	fmt.Printf("%s\n", feed.UpdatedAt)
// 	fmt.Printf("%v\n", feed.LastFetchedAt)
// }
