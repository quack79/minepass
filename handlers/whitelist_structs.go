package handlers

import "regexp"

var minecraftUsername = regexp.MustCompile(`^[A-Za-z0-9_]{3,16}$`)

type WhitelistUserRequest struct {
	Username string `json:"username"`
}

func (r WhitelistUserRequest) IsValid() bool {
	return minecraftUsername.MatchString(r.Username)
}
