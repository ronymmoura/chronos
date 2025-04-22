package http

import (
	"fmt"
	"net/http"

	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/util"

	"github.com/gorilla/mux"
)

type Server struct {
	DB     *db.DB
	Config *util.Config
	Router *mux.Router
}

func NewServer(db *db.DB, config *util.Config, router *mux.Router) *Server {
	return &Server{
		DB:     db,
		Config: config,
		Router: router,
	}
}

func (s *Server) RegisterRoutes() {
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello there!!!")
	})

}
