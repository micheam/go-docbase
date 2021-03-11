package docbase

import "testing"

func Test_buildURL(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty string",
			args{[]string{""}},
			"https://api.docbase.io/",
		},
		{
			"team",
			args{[]string{"team"}},
			"https://api.docbase.io/team",
		},
		{
			"team/domain/posts",
			args{[]string{"team", "domain", "posts"}},
			"https://api.docbase.io/team/domain/posts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildURL(tt.args.paths...); got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
