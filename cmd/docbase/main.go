package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/micheam/go-docbase"
	"github.com/micheam/go-docbase/internal/cli/post"
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
		},
	}
	app.Commands = []*cli.Command{getPost, viewPost, listPosts}
	return app
}

var getPost = &cli.Command{
	Name:      "get",
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
		handler := func(ctx context.Context, post docbase.Post) error {
			tmpl, err := template.New("get-post").Parse(tmplPostDetail)
			if err != nil {
				return err
			}
			err = tmpl.Execute(os.Stdout, post)
			if err != nil {
				return err
			}
			return nil
		}
		return post.Get(c.Context, req, handler)
	},
}

var viewPost = &cli.Command{
	Name:      "view",
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
		handler := func(ctx context.Context, post docbase.Post) error {
			openbrowser(post.URL)
			return nil
		}
		return post.Get(c.Context, req, handler)
	},
}

var listPosts = &cli.Command{
	Name: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
		},
		&cli.IntFlag{
			Name:  "page",
			Value: 1,
		},
		&cli.IntFlag{
			Name:  "per-page",
			Value: 20,
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
		handler := func(ctx context.Context, posts []docbase.Post, meta docbase.Meta) error {
			tmpl, err := template.New("list-posts").Parse(tmplPostsList)
			if err != nil {
				return err
			}
			err = tmpl.Execute(os.Stdout, posts)
			if err != nil {
				return err
			}
			tmpl, err = template.New("meta").Parse(tmplMetaData)
			if err != nil {
				return err
			}
			err = tmpl.Execute(os.Stderr, meta)
			if err != nil {
				return err
			}
			return nil
		}
		return post.List(c.Context, req, handler)
	},
}

var tmplPostDetail = `{{- /* TODO(micheam): Struct2Yaml なライブラリ、あるやろ？感 */ -}}
---
id: {{.ID}} 
title: {{.Title}}
draft: {{.Draft}}
archived: {{.Archived}}
url: {{.URL}}
created_at: {{.CreatedAt}}
updated_at: {{.UpdatedAt}}
---

{{.Body}}
`

var tmplPostsList = `{{range .}}{{printf "%8d" .ID}} {{.Title}}{{"\n"}}{{end}}`

var tmplMetaData = `---
Total: {{.Total}}
{{with .NextPageURL}}Next: {{.}}{{"\n"}}{{end -}} 
{{with .PreviousPageURL}}Prev: {{.}}{{end -}}`

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
