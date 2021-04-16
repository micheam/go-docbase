package docbase

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

// Criteria は、検索時の絞り込みオプションを定義する。
//
// Notes:
//
// 指定可能なキーについては、 docbase 公式のオプションを参照してください。
// https://help.docbase.io/posts/59432?list=%2Fsearch&q=%E6%A4%9C%E7%B4%A2#%E6%A4%9C%E7%B4%A2%E3%82%AA%E3%83%97%E3%82%B7%E3%83%A7%E3%83%B3
//
// なお、Version 1.0 現在、`OR検索` には対応していません。
type Criteria struct {
	Keywords []string
	Include  map[string][]interface{}
	Exclude  map[string][]interface{}
}

func (c *Criteria) MarshalText() (text []byte, err error) {
	sb := new(strings.Builder)

	// - キーワード 検索
	if len(c.Keywords) > 0 {
		_, _ = fmt.Fprintf(sb,
			"%s ", strings.Join(c.Keywords, " "))
	}
	// - AND 検索
	f := func(out io.Writer, format, key string, vals []interface{}) {
		typ := reflect.TypeOf(vals[0])
		if typ.String() == "time.Time" {
			for i := range vals {
				t := vals[i].(time.Time)
				_, _ = fmt.Fprintf(out, format, key, t.Format("2006-01-02"))
			}
			return
		}
		for i := range vals {
			_, _ = fmt.Fprintf(out, format, key, vals[i])
		}
	}
	for key := range c.Include {
		f(sb, "%s:%s ", key, c.Include[key])
	}
	// - NOT 検索
	for key := range c.Exclude {
		f(sb, "-%s:%s ", key, c.Exclude[key])
	}
	// Wrapup
	str := strings.TrimSpace(sb.String())
	return []byte(str), nil
}
