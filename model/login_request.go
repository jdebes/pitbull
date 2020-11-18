package model

import (
	"encoding/json"
	"errors"
)

type LoginRequest struct {
	Username string          `json:"username"`
	Password json.RawMessage `json:"password"`
}

func (r *LoginRequest) Valid() error {
	if len(r.Username) == 0 {
		return errors.New("invalid username")
	}

	if len(r.Password) == 0 || len(r.Password) > 30 {
		return errors.New("invalid password")
	}

	return nil
}
