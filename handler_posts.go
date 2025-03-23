package main

import (
	"context"
	"fmt"
	"log"

	"github.com/szeri323/gator/internal/database"
)

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	var limit int32
	if(len(cmd.args) != 1) {
		limit = 2
	}else {
		fmt.Sscan(cmd.args[0], &limit)
	}
	
	posts, err := s.db.GetUsersPosts(context.Background(), database.GetUsersPostsParams{
	UserID:	user.ID,
	Limit: limit,
})
	if err != nil {
		log.Fatalf("error could not get users posts from db: %v", err)
	}
	for _, post := range posts {
		fmt.Println("@@@@@@@@@")
		fmt.Println(post)
		fmt.Println("@@@@@@@@@")
	}
	return nil
}
