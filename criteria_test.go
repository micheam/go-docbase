package docbase

import (
	"encoding"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var _ encoding.TextMarshaler = (*Criteria)(nil)

func TestCriteria_MarshalText_keywords(t *testing.T) {
	sut := Criteria{[]string{"foo", "bar"}, nil, nil}
	want := []byte("foo bar")
	got, err := sut.MarshalText()
	if err != nil {
		t.Error("Unexpected err: %w", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestCriteria_MarshalText_complex(t *testing.T) {
	sut := Criteria{
		Keywords: []string{"hoge", "moge"},
		Include: map[string][]interface{}{
			"title": {"foo", "bar"},
		},
		Exclude: map[string][]interface{}{
			"liked_by": {"micheam"},
		}}
	want := []byte(
		`hoge moge` +
			` title:foo title:bar` +
			` -liked_by:micheam`)

	got, err := sut.MarshalText()
	if err != nil {
		t.Error("Unexpected err: %w", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestCriteria_MarshalText_Time(t *testing.T) {
	date1, _ := time.Parse("2006-01-02", "2021-01-13")
	date2, _ := time.Parse("2006-01-02", "2051-11-15")
	sut := Criteria{
		Keywords: nil,
		Include: map[string][]interface{}{
			"created_at": {date1},
		},
		Exclude: map[string][]interface{}{
			"updated_at": {date2},
		}}
	want := []byte("created_at:2021-01-13 -updated_at:2051-11-15")
	got, err := sut.MarshalText()
	if err != nil {
		t.Error("Unexpected err: %w", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
