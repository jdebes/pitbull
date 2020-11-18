package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jdebes/pitbull/auth"
	"github.com/jdebes/pitbull/db"
	"github.com/jdebes/pitbull/handler"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type RootHandler struct {
	router *mux.Router
	db     *sqlx.DB
}

func (f *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = db.WithDB(ctx, f.db)

	f.router.ServeHTTP(w, r.WithContext(ctx))
}

func NewServer() *http.Server {
	sqlDb, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	return &http.Server{
		Addr: ":8080",
		Handler: &RootHandler{
			router: newRouter(),
			db:     sqlDb,
		},
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/token", handler.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/sign-up", handler.SignupHandler).Methods(http.MethodPost)

	secured := r.PathPrefix("/auth").Subrouter()
	secured.HandleFunc("/hello", handler.HelloHandler).Methods(http.MethodGet)
	secured.Use(auth.TokenMiddleware)

	r.Handle("/", secured)

	return r
}

func InitLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
