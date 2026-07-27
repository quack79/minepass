package handlers

import (
	"reflect"
	"testing"
)

func TestParseWhitelist(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     []string
	}{
		{
			name:     "players returned by RCON",
			response: "There are 2 whitelisted players: Alex, Builder_01",
			want:     []string{"Alex", "Builder_01"},
		},
		{
			name:     "empty whitelist",
			response: "There are 0 whitelisted players: ",
			want:     []string{},
		},
		{
			name:     "invalid names are ignored",
			response: "There are 2 whitelisted players: Alex, not a player",
			want:     []string{"Alex"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseWhitelist(tt.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseWhitelist() = %v, want %v", got, tt.want)
			}
		})
	}
}
