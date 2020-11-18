package model

type SignupRequest struct {
	LoginRequest
	FirstName string `json:"firstname"`
	Surname   string `json:"surname"`
}

func (r *SignupRequest) Valid() error {
	return nil
}
