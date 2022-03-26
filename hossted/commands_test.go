package hossted

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCommandsMap(t *testing.T) {

	emptyResult := AvailableCommandMap{}

	// Test 1 - Normal Case
	inputA := `
apps:
  - app: demo
    commands: [url, auth]
    values: [example.com, false]
`
	mapA := AvailableCommandMap{}
	mapA["demo.url"] = Command{
		App:     "demo",
		Command: "url",
		Value:   "example.com",
	}
	mapA["demo.auth"] = Command{
		App:     "demo",
		Command: "auth",
		Value:   "false",
	}

	// Test 2 - Mismatched length for commands and values
	inputB := `
apps:
  - app: demo
    commands: [url, auth, aaa]
    values: [example.com, false]
`

	// Test 3 - Invalid yaml. No apps and commands
	inputC := `
app:

`

	// Start test
	type args struct {
		input string
	}
	tests := []struct {
		name               string
		args               args
		want               AvailableCommandMap
		wantErr            bool
		wantErrMsgContains string
	}{
		{
			name: "Normal case",
			args: args{
				input: inputA,
			},
			want: mapA,
		},
		{
			// Error - Mismatched length for command and value
			name: "Mismatched length",
			args: args{
				input: inputB,
			},
			want:               emptyResult,
			wantErr:            true,
			wantErrMsgContains: "does not equal to the length",
		},
		{
			// Error - invalid yaml content
			name: "Invalid yaml",
			args: args{
				input: inputC,
			},
			want:               emptyResult,
			wantErr:            true,
			wantErrMsgContains: "no available apps",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCommandsMap(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCommandsMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommandsMap() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && (tt.wantErrMsgContains != "") {
				assert.Containsf(t, strings.ToLower(err.Error()), strings.ToLower(tt.wantErrMsgContains), "expected error containing %q, got %s", tt.wantErrMsgContains, err)
			}
		})
	}
}
