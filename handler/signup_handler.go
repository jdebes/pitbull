package handler

import (
	"net/http"

	"github.com/jdebes/pitbull/api"
	"github.com/jdebes/pitbull/auth"
	"github.com/jdebes/pitbull/db"
	"github.com/jdebes/pitbull/db/repository"
	"github.com/jdebes/pitbull/model"
	log "github.com/sirupsen/logrus"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var signupRequest model.SignupRequest

	err := api.UnmarshalRequest(w, r, &signupRequest)
	if err != nil {
		log.WithError(err).Error("Unable to unmarshal signup request")
		return
	}

	cryptPassword, err := auth.EncryptPassword(signupRequest.Password)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	dbUser := repository.User{
		FirstName: signupRequest.FirstName,
		Surname:   signupRequest.Surname,
		UserName:  signupRequest.Username,
		Password:  cryptPassword,
	}

	sqlDb, err := db.DB(ctx)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_, err = repository.InsertUser(sqlDb, &dbUser)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
