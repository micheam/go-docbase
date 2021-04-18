package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/micheam/go-docbase"
)

func main() {
	docbase.SetToken(os.Getenv("DOCBASE_TOKEN"))
	var (
		ctx    = context.Background()
		domain = os.Getenv("DOCBASE_DOMAIN")
		param  = url.Values{}
	)
	posts, meta, _ := docbase.ListPosts(ctx, domain, param)
	for i, post := range posts {
		fmt.Printf("%d\t%d\t%s\n", i, post.ID, post.Title)
	}
	fmt.Println("--")
	fmt.Printf("got %d of %d posts\n", len(posts), meta.Total)
}
