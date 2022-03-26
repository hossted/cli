package hossted

import (
	"reflect"
	"testing"
)

func Test_getCommandsMap(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    AvailableCommandMap
		wantErr bool
	}{
		// TODO: Add test cases.
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
