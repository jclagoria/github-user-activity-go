package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantUser string
		wantType string
		wantErr  bool
	}{
		{
			name:     "username only",
			args:     []string{"github-activity", "jclagoria"},
			wantUser: "jclagoria",
			wantType: "",
		},
		{
			name:     "username with --type",
			args:     []string{"github-activity", "jclagoria", "--type", "PushEvent"},
			wantUser: "jclagoria",
			wantType: "PushEvent",
		},
		{
			name:    "no arguments",
			args:    []string{"github-activity"},
			wantErr: true,
		},
		{
			name:    "--type without value",
			args:    []string{"github-activity", "jclagoria", "--type"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"github-activity", "jclagoria", "--unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, typ, err := ParseArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseArgs() error = nil, wantErr = true")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseArgs() unexpected error: %v", err)
				return
			}
			if user != tt.wantUser {
				t.Errorf("ParseArgs() user = %q, want %q", user, tt.wantUser)
			}
			if typ != tt.wantType {
				t.Errorf("ParseArgs() type = %q, want %q", typ, tt.wantType)
			}
		})
	}
}
