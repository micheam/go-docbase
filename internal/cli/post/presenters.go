package post

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/micheam/go-docbase"
)

func WritePostToConsole(ctx context.Context, post docbase.Post) error {
	// TODO(micheam): Struct2Yaml なライブラリを探して切り替え
	const tmplPostDetail = `---
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

func OpenBrowser(_ context.Context, post docbase.Post) error {
	openbrowser(post.URL)
	return nil
}

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

func BuildListResultHandler(withMeta bool) (ListResonseHandler, error) {
	const _tmplPostsList = `{{range .}}{{printf "%d\t%s" .ID (summary .)}}{{"\n"}}{{end}}`
	tmplPostsList, err := template.New("list-posts").Funcs(template.FuncMap{
		"summary": summary,
	}).Parse(_tmplPostsList)
	if err != nil {
		return nil, err
	}
	const _tmplMetaData = `---
Total: {{.Total}}
{{with .NextPageURL}}Next: {{.}}{{"\n"}}{{end -}} 
{{with .PreviousPageURL}}Prev: {{.}}{{end -}}`
	tmplMetaData, err := template.New("meta").Parse(_tmplMetaData)
	if err != nil {
		return nil, err
	}
	if withMeta {
		return func(ctx context.Context, posts []docbase.Post, meta docbase.Meta) error {
			err := tmplPostsList.Execute(os.Stdout, posts)
			if err != nil {
				return err
			}
			return tmplMetaData.Execute(os.Stdout, meta)
		}, nil
	}
	return func(ctx context.Context, posts []docbase.Post, _ docbase.Meta) error {
		return tmplPostsList.Execute(os.Stdout, posts)
	}, nil
}
