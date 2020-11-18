package auth

import (
	"net/http"

	"github.com/jdebes/pitbull/api"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func CompareLoginPassword(w http.ResponseWriter, password []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		api.Debug(w, err, http.StatusUnauthorized)
	}

	return err == nil
}
