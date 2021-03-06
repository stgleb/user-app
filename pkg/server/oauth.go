package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"user-app/pkg/user"
)

var (
	googleOauthConfig *oauth2.Config
)

func InitOAuth(clientId, clientSecret string) {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/contacts",
			"https://www.googleapis.com/auth/contacts.readonly",
			"https://www.googleapis.com/auth/directory.readonly",
			"https://www.googleapis.com/auth/user.addresses.read",
			"https://www.googleapis.com/auth/user.birthday.read",
			"https://www.googleapis.com/auth/user.emails.read",
			"https://www.googleapis.com/auth/user.gender.read",
			"https://www.googleapis.com/auth/user.organization.read",
			"https://www.googleapis.com/auth/user.phonenumbers.read",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (s *Server) loginGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(stateToken)
	http.Redirect(w, r, fmt.Sprintf("%s&&type=login", url), http.StatusFound)
}

func (s *Server) signUpGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(stateToken)
	http.Redirect(w, r, fmt.Sprintf("%s&&type=signup", url), http.StatusFound)
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	u, err := getUserInfo(r.Context(), state, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Codepath for login of existing user
	_, err = s.repo.FindById(r.Context(), u.Id)
	if err == user.NotFound {
		_, err := s.repo.Store(r.Context(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Store user info in session cookie
	if err := s.loginUser(u.Id, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/user/%s", u.Id), http.StatusFound)
}

func getUserInfo(ctx context.Context, state string, code string) (u *user.User, err error) {
	if state != stateToken {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer func() {
		err = response.Body.Close()
	}()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	u = &user.User{}
	if err := json.Unmarshal(contents, u); err != nil {
		return nil, err
	}
	return u, nil
}
