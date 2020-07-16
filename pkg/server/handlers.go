package server

import (
	"net/http"
)

var (
	// TODO(stgleb): randomize it
	stateToken string = "state"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) userInfo(w http.ResponseWriter, r *http.Request) {
}


func (s *Server) editUserInfo(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) editInfo(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) loginGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(stateToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
