package hossted

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetDefaultApp(t *testing.T) {
	type fields struct {
		Email        string
		UserToken    string
		SessionToken string
		EndPoint     string
		UUIDPath     string
		Applications []ConfigApplication
	}
	type args struct {
		pwd string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               string
		wantErr            bool
		wantErrMsgContains string
	}{
		{
			name: "Test current directory",
			fields: fields{
				Email:        "",
				UserToken:    "",
				SessionToken: "",
				EndPoint:     "",
				UUIDPath:     "",
				Applications: []ConfigApplication{},
			},
			args:    args{pwd: "/opt/gitbucket"},
			want:    "gitbucket",
			wantErr: false,
		},
		{
			name: "Test default app",
			fields: fields{
				Email:        "",
				UserToken:    "",
				SessionToken: "",
				EndPoint:     "",
				UUIDPath:     "",
				Applications: []ConfigApplication{
					ConfigApplication{
						AppName: "prometheus",
						AppPath: "",
					},
				},
			},
			args:    args{pwd: "/sdf"},
			want:    "prometheus",
			wantErr: false,
		},
		{
			name: "No default apps - multiple apps in config",
			fields: fields{
				Email:        "",
				UserToken:    "",
				SessionToken: "",
				EndPoint:     "",
				UUIDPath:     "",
				Applications: []ConfigApplication{
					ConfigApplication{
						AppName: "prometheus",
						AppPath: "",
					},
					ConfigApplication{
						AppName: "gitbucket",
						AppPath: "",
					},
				},
			},
			args:               args{pwd: "/sdf"},
			want:               "",
			wantErr:            true,
			wantErrMsgContains: "no default apps", // As contains multiple
		},
		{
			name: "Default folders, with multiple apps in config",
			fields: fields{
				Email:        "",
				UserToken:    "",
				SessionToken: "",
				EndPoint:     "",
				UUIDPath:     "",
				Applications: []ConfigApplication{
					ConfigApplication{
						AppName: "prometheus",
						AppPath: "",
					},
					ConfigApplication{
						AppName: "gitbucket",
						AppPath: "",
					},
				},
			},
			args:    args{pwd: "/opt/wordpress"},
			want:    "wordpress",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Email:        tt.fields.Email,
				UserToken:    tt.fields.UserToken,
				SessionToken: tt.fields.SessionToken,
				EndPoint:     tt.fields.EndPoint,
				UUIDPath:     tt.fields.UUIDPath,
				Applications: tt.fields.Applications,
			}
			got, err := c.GetDefaultApp(tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.GetDefaultApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Config.GetDefaultApp() = %v, want %v", got, tt.want)
			}
			if tt.wantErr && (tt.wantErrMsgContains != "") {
				assert.Containsf(t, strings.ToLower(err.Error()), strings.ToLower(tt.wantErrMsgContains), "expected error containing %q, got %s", tt.wantErrMsgContains, err)
			}
		})
	}
}
