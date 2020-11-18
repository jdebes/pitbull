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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest model.LoginRequest
	var user repository.User
	ctx := r.Context()

	err := api.UnmarshalRequest(w, r, &loginRequest)
	if err != nil {
		log.WithError(err).Error("Unable to unmarshal log request")
		return
	}

	sqlDb, err := db.DB(ctx)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	err = repository.UserByUserName(sqlDb, loginRequest.Username, &user)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	if auth.CompareLoginPassword(w, loginRequest.Password, user.Password) {
		jwt, err := auth.BuildJwt(loginRequest.Username)
		if err != nil {
			api.Error(w, err, http.StatusBadRequest)
			return
		}

		err = api.MarshalResponse(w, map[string]string{"token": jwt})
		if err != nil {
			log.WithError(err).Error("Unable to marshal login response")
		}
	}
}
