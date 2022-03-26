package hossted

import (
	"reflect"
	"testing"
)

func Test_getCommandsMap(t *testing.T) {
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

	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    AvailableCommandMap
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				input: inputA,
			},
			want: mapA,
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
		})
	}
}
