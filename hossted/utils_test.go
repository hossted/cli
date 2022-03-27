package hossted

import "testing"

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
