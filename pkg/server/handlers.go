package server

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"user-app/pkg/user"
)

var (
	// TODO(stgleb): randomize it
	stateToken string = "state"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect unauthenticated user to login page
	userId := session.Values["user"]
	if userId == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/user/%v", userId), http.StatusFound)
}

func (s *Server) userInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	u, err := s.repo.FindById(r.Context(), userId)
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
	u, err := s.repo.FindById(r.Context(), userId)
	if err != nil && err == user.NotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		t := s.templateMap[UserEditInfo]
		t.Execute(w, u)
	} else {
		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		address := r.FormValue("address")
		u, err := s.repo.FindByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.Name = name
		u.Telephone = phone
		u.Address = address
		u.Email = email
		err = s.repo.Update(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/user/%s", u.Id), http.StatusFound)
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[Login]
		t.Execute(w, nil)
	} else {
		password := r.FormValue("password")
		email := r.FormValue("email")
		u, err := s.repo.FindByEmail(r.Context(), email)
		if err == user.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		passwordHash := sha256.Sum256([]byte(password))
		hexPassword := hex.EncodeToString(passwordHash[:])
		if !strings.EqualFold(u.PasswordHash, hexPassword) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err := s.loginUser(u.Id, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/user/%s", u.Id), http.StatusFound)
	}
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[SignUp]
		t.Execute(w, struct {
			ApiKey string
		}{
			ApiKey: s.googleApiKey,
		})
	} else {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		phone := r.FormValue("phone")
		address := r.FormValue("address")
		if u, _ := s.repo.FindByEmail(r.Context(), email); u != nil {
			http.Error(w, fmt.Sprintf("user with email %s already exists", email),
				http.StatusConflict)
			return
		}
		passwordHash := sha256.Sum256([]byte(password))
		hexPassword := hex.EncodeToString(passwordHash[:])
		u := &user.User{
			Id:           uuid.New().String(),
			Name:         username,
			Email:        email,
			PasswordHash: hexPassword,
			Telephone:    phone,
			Address:      address,
		}
		userId, err := s.repo.Store(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Store user info in session cookie
		if err := s.loginUser(u.Id, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/user/%s", userId), http.StatusFound)
	}
}

func (s *Server) forgetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t := s.templateMap[ForgotPassword]
		t.Execute(w, nil)
		return
	}
	email := r.FormValue("email")
	_, err := s.repo.FindByEmail(r.Context(), email)
	if err == user.NotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := s.sendResetPasswordEmail(r.Context(), r.Host, email); err != nil {
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
		token, err := s.repo.GetByToken(r.Context(), values[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := s.repo.DisableToken(r.Context(), token.Value); err != nil {
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
		u, err := s.repo.FindByEmail(r.Context(), email)
		if err == user.NotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash := sha256.Sum256([]byte(password))
		u.PasswordHash = hex.EncodeToString(hash[:])
		if err := s.repo.Update(r.Context(), u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// logout revokes authentication for a user
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user"] = nil
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) loginUser(userId string, w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "user")
	if err != nil {
		return errors.New("error getting user")
	}
	session.Values["user"] = userId
	err = session.Save(r, w)
	if err != nil {
		return errors.New("error saving user")
	}
	return nil
}
