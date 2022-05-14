package hossted

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBool(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name               string
		args               args
		want               bool
		wantErr            bool
		wantErrMsgContains string
	}{
		{
			name: "Input as true",
			args: args{
				in: "true",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Input as false",
			args: args{
				in: "false",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Input as something else",
			args: args{
				in: "abc",
			},
			want:               false,
			wantErr:            true,
			wantErrMsgContains: "Only true/false is supported",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertBool(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertBool() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && (tt.wantErrMsgContains != "") {
				assert.Containsf(t, strings.ToLower(err.Error()), strings.ToLower(tt.wantErrMsgContains), "expected error containing %q, got %s", tt.wantErrMsgContains, err)
			}
		})
	}
}
