package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"user-app/pkg/user"
	"user-app/pkg/user/memory"
)

type Server struct {
	repo user.Repository
	srv  *http.Server
}

func NewServer(addr string) *Server {
	router := mux.NewRouter()
	srv := &Server{}
	router.HandleFunc("/", srv.indexHandler)
	router.HandleFunc("/:id", srv.userInfo).Methods(http.MethodGet)
	router.HandleFunc("/:id", srv.editInfo).Methods(http.MethodPut)
	router.HandleFunc("/login", srv.login)
	router.HandleFunc("/signup", srv.signUp)
	router.HandleFunc("/login/google", srv.loginGoogle)
	router.HandleFunc("/signup/google", srv.signUpGoogle)
	router.HandleFunc("/callback", srv.callback)

	httpSrv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	srv.srv = httpSrv
	srv.repo = memory.NewRepository()

	return srv
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}
