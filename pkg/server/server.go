package server

import (
	"html/template"
	"net/http"
	"user-app/pkg/user/repository"

	"github.com/gorilla/mux"

	"user-app/pkg/user/repository/memory"
)

type Server struct {
	SmtpServerHost string
	SmtpServerPort int
	SmtpUser       string
	SmtpPassowrd   string
	googleApiKey   string
	repo           repository.Repository
	srv            *http.Server
	templateMap    map[string]*template.Template
}

func NewServer(addr, templatesDir, smtpServerHost string, smtpServerPort int, smtpUser, smtpPassword, googleApiKey string) (*Server, error) {
	router := mux.NewRouter()
	srv := &Server{
		SmtpServerHost: smtpServerHost,
		SmtpServerPort: smtpServerPort,
		SmtpUser:       smtpUser,
		SmtpPassowrd:   smtpPassword,
		googleApiKey:   googleApiKey,
		templateMap:    make(map[string]*template.Template),
	}
	router.HandleFunc("/user/{id}", srv.userInfo).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}/edit", srv.editUserInfo).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/login", srv.login)
	router.HandleFunc("/signup", srv.signUp)
	router.HandleFunc("/login/google", srv.loginGoogle)
	router.HandleFunc("/signup/google", srv.signUpGoogle)
	router.HandleFunc("/callback", srv.callback)
	router.HandleFunc("/password/forget", srv.forgetPassword).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/password/reset", srv.resetPassword).Methods(http.MethodGet, http.MethodPost)
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
