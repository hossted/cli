package hossted

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				format: "domain",
			},
			want: true,
		},
		{
			name: "Simple test two",
			args: args{
				in:     "https://www.example.com",
				format: "domain",
			},
			want: true,
		},
		{
			name: "Simple test three",
			args: args{
				in:     "http://www.example.com",
				format: "domain",
			},
			want: true,
		},
		{
			name: "Simple test four",
			args: args{
				in:     "https://www.example.com/",
				format: "domain",
			},
			want: true,
		},
		{
			name: "Simple test five",
			args: args{
				in:     "abccc",
				format: "domain",
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

func Test_writeProtected(t *testing.T) {
	type args struct {
		filepath string
		b        []byte
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		wantErrMsgContains string
	}{
		{
			name: "File does not exist",
			args: args{
				filepath: "/tmp/ddddeeeeddddd.txt",
				b:        []byte("abcd"),
			},
			wantErr:            true,
			wantErrMsgContains: "protected file does not exist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeProtected(tt.args.filepath, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("writeProtected() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && (tt.wantErrMsgContains != "") {
				assert.Containsf(t, strings.ToLower(err.Error()), strings.ToLower(tt.wantErrMsgContains), "expected error containing %q, got %s", tt.wantErrMsgContains, err)
			}
		})
	}
}
