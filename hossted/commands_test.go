package hossted

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCommandsMap(t *testing.T) {

	emptyResult := AvailableCommandMap{}

	generalCmdEmpty := ""
	generalCmd := ""
	_ = generalCmd

	// Test 1 - Normal Case
	appCmdA := `
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
	appCmdB := `
apps:
  - app: demo
    commands: [url, auth, aaa]
    values: [example.com, false]
`

	// Test 3 - Invalid yaml. No apps and commands
	appCmdC := `
app:

`

	// Start test
	type args struct {
		generalCmd string
		appCmd     string
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
				generalCmd: generalCmdEmpty,
				appCmd:     appCmdA,
			},
			want: mapA,
		},
		{
			// Error - Mismatched length for command and value
			name: "Mismatched length",
			args: args{
				generalCmd: generalCmdEmpty,
				appCmd:     appCmdB,
			},
			want:               emptyResult,
			wantErr:            true,
			wantErrMsgContains: "does not equal to the length",
		},
		{
			// Error - invalid yaml content
			name: "Invalid yaml",
			args: args{
				generalCmd: generalCmdEmpty,
				appCmd:     appCmdC,
			},
			want:               emptyResult,
			wantErr:            true,
			wantErrMsgContains: "no available apps",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCommandsMap(tt.args.generalCmd, tt.args.appCmd)
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

func TestCheckCommands(t *testing.T) {
	type args struct {
		app     string
		command string
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		wantErrMsgContains string
	}{
		{
			// App Command not in predefined list.
			// Assume there are some commands defined
			name: "App Command not predefined",
			args: args{
				app:     "aaa",
				command: "bbb",
			},
			wantErr:            true,
			wantErrMsgContains: "is not supported",
		},
		{
			// App Command not in predefined list.
			// Assume there are some commands defined
			name: "General command - remote-support should be avilable",
			args: args{
				app:     "general",
				command: "remote-support",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCommands(tt.args.app, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckCommands() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && (tt.wantErrMsgContains != "") {
				assert.Containsf(t, strings.ToLower(err.Error()), strings.ToLower(tt.wantErrMsgContains), "expected error containing %q, got %s", tt.wantErrMsgContains, err)
			}
		})
	}
}
