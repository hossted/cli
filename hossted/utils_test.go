package hossted

import (
	"testing"
)

func Test_verifyInputFormat(t *testing.T) {
	type args struct {
		in     string
		format string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Simple test one",
			args: args{
				in:     "example.com",
				format: "url",
			},
			want: true,
		},
		{
			name: "Simple test two",
			args: args{
				in:     "https://www.example.com",
				format: "url",
			},
			want: true,
		},
		{
			name: "Simple test three",
			args: args{
				in:     "http://www.example.com",
				format: "url",
			},
			want: true,
		},
		{
			name: "Simple test four",
			args: args{
				in:     "https://www.example.com/",
				format: "url",
			},
			want: true,
		},
		{
			name: "Simple test five",
			args: args{
				in:     "abccc",
				format: "url",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyInputFormat(tt.args.in, tt.args.format); got != tt.want {
				t.Errorf("verifyInputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAppFilePath(t *testing.T) {
	type args struct {
		base     string
		relative string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			// Simple test - assume /usr/local/bin exists
			name: "Simple test",
			args: args{
				base:     "/usr/local",
				relative: "bin",
			},
			want: "/usr/local/bin",
		},
		{
			// Simple test - assume /usr/local/bin exists
			name: "Simple test 2 ",
			args: args{
				base:     "/usr/local/",
				relative: "bin",
			},
			want: "/usr/local/bin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAppFilePath(tt.args.base, tt.args.relative)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAppFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAppFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
