package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"user-app/pkg/user"
)

var (
	// TODO(stgleb): randomize it
	stateToken string = "state"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) userInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	u, err := s.repo.FindById(userId)
	if err != nil && err == user.NotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := s.templateMap[UserInfo]
	t.Execute(w, u)
}

func (s *Server) editUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	u, err := s.repo.FindById(userId)
	if err != nil && err == user.NotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		t := s.templateMap[UserEditInfo]
		t.Execute(w, u)
	} else {

	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	t := s.templateMap[Login]
	t.Execute(w, nil)
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	t := s.templateMap[SignUp]
	t.Execute(w, nil)
}

func (s *Server) forgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[ForgotPassword]
		t.Execute(w, nil)
		return
	}

	email := r.FormValue("email")
	log.Println(email)
}