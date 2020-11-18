package main

import (
	"github.com/jdebes/pitbull/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	server.InitLogging()
	s := server.NewServer()

	log.WithField("address", s.Addr).Info("Starting up app")
	err := s.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
