package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/micheam/go-docbase"
	"github.com/micheam/go-docbase/internal/cli/post"
	"github.com/micheam/go-docbase/internal/cli/tag"
	"github.com/micheam/go-docbase/internal/pointer"
	"github.com/urfave/cli/v2"
)

var version = "0.1.0"

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(io.Discard)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "docbase"
	app.Usage = "CLI Client-Tool for DocBase API"
	app.Version = version
	app.Authors = []*cli.Author{
		{
			Name:  "Michito Maeda",
			Email: "michito.maeda@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"vv"},
			EnvVars: []string{"DOCBASE_VERBOSE", "DOCBASE_DEBUG", "DEBUG"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "`ACCESS_TOKEN` for docbase API",
			EnvVars: []string{"DOCBASE_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "domain",
			EnvVars: []string{"DOCBASE_DOMAIN"},
			Usage:   "`NAME` on docbase.io",
		},
	}
	app.Commands = []*cli.Command{
		getPost, viewPost, listPosts, createPost,
		tags,
	}
	return app
}

var getPost = &cli.Command{
	Name:      "get",
	Usage:     "Get post content on docbase.io",
	ArgsUsage: "POST_ID",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		postID, err := docbase.ParsePostID(c.Args().First())
		if err != nil {
			return err
		}
		req := post.GetRequest{
			Domain: c.String("domain"),
			ID:     postID,
		}
		return post.Get(c.Context, req, post.WritePostToConsole)
	},
}

var viewPost = &cli.Command{
	Name:      "view",
	Usage:     "View post on docbase.io",
	ArgsUsage: "POST_ID",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		postID, err := docbase.ParsePostID(c.Args().First())
		if err != nil {
			return err
		}
		req := post.GetRequest{
			Domain: c.String("domain"),
			ID:     postID,
		}
		return post.Get(c.Context, req, post.OpenBrowser)
	},
}

var listPosts = &cli.Command{
	Name:  "list",
	Usage: "Search and list posts on docbase.io",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "`options` to narrow down the search. ex: groups,contributors, etc.",
		},
		&cli.IntFlag{
			Name:    "page",
			Aliases: []string{"p"},
			Value:   1,
			Usage:   "`num` of posts on a page",
		},
		&cli.IntFlag{
			Name:    "per-page",
			Aliases: []string{"pp"},
			Value:   20,
			Usage:   "`num` of page",
		},
		&cli.BoolFlag{
			Name:    "meta",
			Aliases: []string{"m"},
			Usage:   "Display META-Fields (Total,Previous,Next) on footer",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := post.ListRequest{Domain: c.String("domain")}
		if c.String("query") != "" {
			req.Query = pointer.StringPtr(c.String("query"))
		}
		if c.Int("page") != 0 {
			req.Page = pointer.IntPtr(c.Int("page"))
		}
		if c.Int("per-page") != 0 {
			req.PerPage = pointer.IntPtr(c.Int("per-page"))
		}
		presenter, err := post.BuildListResultHandler(c.Bool("meta"))
		if err != nil {
			return err
		}
		return post.List(c.Context, req, presenter)
	},
}

// TODO(micheam): 設定ファイルで指定可能にする
var defaultTitle = func() string {
	now := time.Now()
	return fmt.Sprintf("%s 作業メモ", now.Format("2006-01-02"))
}

var createPost = &cli.Command{
	Name:      "new",
	Usage:     "Create new post.",
	ArgsUsage: "title",
	Flags: []cli.Flag{
		// TODO(micheam): option `--dradt`
		// TODO(micheam): option `--notice`
		// TODO(micheam): option `--tags`
		// TODO(micheam): option `--scope`
		// TODO(micheam): option `--groups`
		&cli.StringFlag{
			Name:     "title",
			Aliases:  []string{"t"},
			Usage:    "`VALUE` of title",
			Required: false,
			Value:    defaultTitle(),
		},
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"f"},
			Usage:    "`PATH` of input file",
			Required: true, // TODO(micheam): change to OPTIONAL
			//   省略された場合はエディタを起動するようにしたい
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := post.CreateRequest{
			Title:  c.String("title"),
			Domain: c.String("domain"),
		}
		if len(c.String("file")) != 0 {
			filepath := c.String("file")
			file, err := os.Open(filepath)
			if err != nil {
				return fmt.Errorf("cant open %q: %w", filepath, err)
			}
			defer func() { _ = file.Close() }()
			req.Body = file
		}
		presenter := func(ctx context.Context, post *docbase.Post) error {
			fmt.Println(post.URL)
			return nil
		}
		return post.Create(c.Context, req, presenter)
	},
}

var tags = &cli.Command{
	Name:  "tags",
	Usage: "Show tags of group",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		req := tag.ListRequest{
			Domain: c.String("domain"),
		}
		presenter := func(ctx context.Context, tags []docbase.Tag) error {
			for _, tag := range tags {
				fmt.Println(tag.Name)
			}
			return nil
		}
		return tag.List(c.Context, req, presenter)
	},
}
