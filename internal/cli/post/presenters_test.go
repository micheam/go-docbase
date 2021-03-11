package post

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_marshal(t *testing.T) {
	data := struct {
		F1 string   `yaml:"title"`
		F2 int      `yaml:"id"`
		F3 []string `yaml:"tags"`
		F4 string   `yaml:"-"`
		F5 string   `yaml:"foo,omitempty"`
	}{
		F1: "Hello, Yaml World.",
		F2: 33224455,
		F3: []string{"dev", "golang"},
	}
	want := `title: Hello, Yaml World.
id: 33224455
tags:
- dev
- golang
`
	got := marshal(data)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):%s\n", diff)
	}
}
