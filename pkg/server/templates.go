package server

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

const (
	Login          = "login"
	SignUp         = "signup"
	UserInfo       = "user"
	UserEditInfo   = "user_edit"
	ForgotPassword = "forgot_password"
	Email          = "email"
	Reset          = "reset"
)

// ParseTemplates parse templates in template folder
func (s *Server) ParseTemplates(templateMap map[string]*template.Template, dirname string) error {
	if len(dirname) == 0 {
		return nil
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() {
			fullName := path.Join(dirname, f.Name())
			f, err := os.Open(fullName)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			lastTerm := len(strings.Split(f.Name(), "/"))
			name := strings.Split(strings.Split(f.Name(), "/")[lastTerm-1], ".")[0]
			t, err := template.New(name).Parse(string(data))
			if err != nil {
				return errors.Wrapf(err, "failed to parse %s template", fullName)
			}
			templateMap[name] = t
		}
	}

	return nil
}
