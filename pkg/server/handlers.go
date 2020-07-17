package server

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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
	  // TODO(stgleb): Implement this
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[Login]
		t.Execute(w, nil)
	} else {
		password := r.FormValue("password")
		email := r.FormValue("email")

		u, err := s.repo.FindByEmail(email)
		if err == user.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		passwordHash := sha256.Sum256([]byte(password))
		if !bytes.Equal(u.PasswordHash, passwordHash[:]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/user/%s", u.Id), http.StatusMovedPermanently)
	}
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	t := s.templateMap[SignUp]
	t.Execute(w, nil)
}

func (s *Server) forgetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[ForgotPassword]
		t.Execute(w, nil)
		return
	}

	email := r.FormValue("email")
	_, err := s.repo.FindByEmail(email)
	if err == user.NotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.sendResetPasswordEmail(r.Host, email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("Reset email has been sent to %s", email)))
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		values, ok := r.URL.Query()["token"]
		if !ok || len(values[0]) < 1 {
			http.Error(w, "token is missing", http.StatusBadRequest)
			return
		}
		token, err := s.repo.GetByToken(values[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token.Used = true
		if err := s.repo.StoreToken(token);err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t := s.templateMap[Reset]
		t.Execute(w, struct {
			Email string
		}{
			Email: token.Email,
		})
	} else {
		password := r.FormValue("password")
		email := r.FormValue("email")
		u, err := s.repo.FindByEmail(email)
		if err == user.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash := sha256.Sum256([]byte(password))
		u.PasswordHash = hash[:]
		if _, err := s.repo.Store(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Password has been changed"))
	}
}
