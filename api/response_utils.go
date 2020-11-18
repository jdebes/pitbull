package api

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func MarshalResponse(w http.ResponseWriter, model interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(model)
	if err != nil {
		return err
	}

	return nil
}

func Error(w http.ResponseWriter, raisedError error, status int) {
	writeErrResponseWithLogLevel(w, raisedError, http.StatusText(status), status, log.ErrorLevel)
}

func Info(w http.ResponseWriter, raisedError error, status int) {
	writeErrResponseWithLogLevel(w, raisedError, http.StatusText(status), status, log.InfoLevel)
}

func Debug(w http.ResponseWriter, raisedError error, status int) {
	writeErrResponseWithLogLevel(w, raisedError, http.StatusText(status), status, log.DebugLevel)
}

func writeErrResponseWithLogLevel(w http.ResponseWriter, raisedError error, message string, status int, level log.Level) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	jsonResponse, err := json.Marshal(map[string]string{"error": message})
	if err != nil {
		log.WithError(err).Error("Failed to serialise response error")
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.WithError(err).Error("Failed to write error response")
	}

	if raisedError != nil {
		log.WithError(raisedError).Log(level, message)
	} else {
		log.WithTime(time.Now()).Log(level, message)
	}
}
