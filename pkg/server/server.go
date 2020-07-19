package server

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"

	"user-app/pkg/user/repository"
)

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)
	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

type Server struct {
	SmtpServerHost string
	SmtpServerPort int
	SmtpUser       string
	SmtpPassword   string
	EmailFrom      string
	googleApiKey   string
	repo           repository.Repository
	srv            *http.Server
	templateMap    map[string]*template.Template
}

// NewServer creates configured server instance ready to run
func NewServer(addr, templatesDir, repositoryType string, repositoryOpts repository.Opts, smtpServerHost string,
	smtpServerPort int, smtpUser, smtpPassword, googleApiKey string) (*Server, error) {
	router := mux.NewRouter()
	srv := &Server{
		SmtpServerHost: smtpServerHost,
		SmtpServerPort: smtpServerPort,
		SmtpUser:       smtpUser,
		SmtpPassword:   smtpPassword,
		googleApiKey:   googleApiKey,
		templateMap:    make(map[string]*template.Template),
	}
	router.HandleFunc("/", srv.indexHandler)
	router.HandleFunc("/login", srv.login)
	router.HandleFunc("/signup", srv.signUp)
	router.HandleFunc("/login/google", srv.loginGoogle)
	router.HandleFunc("/signup/google", srv.signUpGoogle)
	router.HandleFunc("/callback", srv.callback)
	router.HandleFunc("/password/forget", srv.forgetPassword).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/password/reset", srv.resetPassword).Methods(http.MethodGet, http.MethodPost)

	router.HandleFunc("/user/{id}", secureMiddleware(srv.userInfo)).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}/edit", secureMiddleware(srv.editUserInfo)).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/logout", secureMiddleware(srv.logout))
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	srv.srv = httpSrv
	repo, err := NewRepository(repositoryType, repositoryOpts)
	if err != nil {
		return nil, err
	}
	srv.repo = repo
	err = srv.ParseTemplates(srv.templateMap, templatesDir)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}
