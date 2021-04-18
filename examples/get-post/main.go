package main

import (
	"context"
	"fmt"
	"os"

	"github.com/micheam/go-docbase"
)

func main() {
	docbase.SetToken(os.Getenv("DOCBASE_TOKEN"))
	var (
		ctx    = context.Background()
		domain = os.Getenv("DOCBASE_DOMAIN")
		postID = docbase.PostID(1863830) // 記事ID
	)
	post, _ := docbase.GetPost(ctx, domain, postID)
	fmt.Printf("%d: %s\n", post.ID, post.Title)
	fmt.Println(post.Body)
}
