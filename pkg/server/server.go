package server

import (
	"html/template"
	"net/http"
	"user-app/pkg/user/repository"

	"github.com/gorilla/mux"

	"user-app/pkg/user/repository/memory"
)

type Server struct {
	repo repository.Repository
	srv  *http.Server
	templateMap map[string]*template.Template
}

func NewServer(addr, templatesDir string) (*Server, error) {
	router := mux.NewRouter()
	srv := &Server{
		templateMap: make(map[string]*template.Template),
	}
	router.HandleFunc("/user/{id}", srv.userInfo).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}/edit", srv.editUserInfo).Methods(http.MethodPut, http.MethodGet)
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
	err := srv.ParseTemplates(srv.templateMap, templatesDir)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}
