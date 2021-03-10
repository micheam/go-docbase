package post

import (
	"context"
	"log"

	"github.com/micheam/go-docbase"
)

type GetRequest struct {
	ID     docbase.PostID
	Domain string
}

type GetResonseHandler func(ctx context.Context, post docbase.Post) error

func Get(ctx context.Context, req GetRequest, handler GetResonseHandler) error {
	log.Printf("get post with req: %v", req)
	post, err := docbase.GetPost(ctx, req.Domain, req.ID)
	if err != nil {
		return err
	}
	return handler(ctx, *post)
}
