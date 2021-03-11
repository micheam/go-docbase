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
	"github.com/micheam/go-docbase/internal/text"
	"gopkg.in/yaml.v2"
)

type PostDetailViewData struct {
	ID         int           `yaml:"id"`
	Title      string        `yaml:"title"`
	Body       string        `yaml:"-"`
	Draft      bool          `yaml:"draft"`
	Archived   bool          `yaml:"archived"`
	URL        string        `yaml:"url"`
	CreatedAt  string        `yaml:"created_at"` // ISO 8601
	UpdatedAt  string        `yaml:"updated_at"` // ISO 8601
	Scope      string        `yaml:"scope"`
	SharingURL string        `yaml:"-"`
	Tags       []string      `yaml:"tags"`
	UserName   string        `yaml:"user"`
	Stars      int           `yaml:"stars_count"`
	GoodJob    int           `yaml:"good_jobs_count"`
	Comments   []interface{} `yaml:"-"`
	Groups     []interface{} `yaml:"-"`
}

func NewPostDetailViewData(src docbase.Post) *PostDetailViewData {
	tags := []string{}
	for i := range src.Tags {
		tags = append(tags, src.Tags[i].Name)
	}
	return &PostDetailViewData{
		ID:         src.ID.Int(),
		Title:      src.Title,
		Body:       src.Body,
		Draft:      src.Draft,
		Archived:   src.Archived,
		URL:        src.URL,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		Scope:      string(src.Scope),
		SharingURL: src.SharingURL,
		Tags:       tags,
		UserName:   src.User.Name,
		Stars:      src.Stars,
		GoodJob:    src.GoodJob,
		Comments:   src.Comments,
		Groups:     src.Groups,
	}
}

func marshal(v interface{}) string {
	a, _ := yaml.Marshal(v)
	return string(a)
}

func WritePostToConsole(ctx context.Context, post docbase.Post) error {
	// TODO(micheam): Win対応
	//   現状、Body の改行コードを一律変換してしまっている。
	const tmplPostDetail = `---
{{marshal .}}
---

{{dos2unix .Body}}
`
	tmpl, err := template.New("get-post").Funcs(
		template.FuncMap{
			"marshal":  marshal,
			"dos2unix": text.Dos2Unix,
		}).Parse(tmplPostDetail)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, NewPostDetailViewData(post))
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
