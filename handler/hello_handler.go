package handler

import (
	"net/http"

	"github.com/jdebes/pitbull/api"
	"github.com/jdebes/pitbull/auth"
)

const htmlResponse = "<!DOCTYPE html>\n<html>\n<head>\n<title>Page Title</title>\n</head>\n<body>\n\n<h1>This is a Heading</h1>\n<p>This is a paragraph.</p>\n\n</body>\n</html>"

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ContextUserClaims(r.Context())
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
	}

	w.Header().Set("X-Username", claims.Username)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponse))
}
