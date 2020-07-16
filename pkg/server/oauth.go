package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func (s *Server) signUpGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(stateToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	if state != stateToken {
		http.Error(w, "invalid oauth state", http.StatusInternalServerError)
		return
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("code exchange failed: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed getting user info: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed reading response body: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Content: %s\n", contents)
}
