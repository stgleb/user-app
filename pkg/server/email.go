package server

import (
	"bytes"
	"fmt"
	"net/smtp"

	"github.com/google/uuid"

	"user-app/pkg/user"
)

func (s *Server) sendResetPasswordEmail(hostName string, email string) error {
	t := s.templateMap[Email]
	buffer := &bytes.Buffer{}
	tokenValue := uuid.New().String()
	token := &user.Token{
		Email: email,
		Value: tokenValue,
		Used:  false,
	}
	if err := s.repo.StoreToken(token); err != nil {
		return err
	}
	err := t.Execute(buffer, struct {
		Email    string
		HostName string
		Token    string
	}{
		Email:    email,
		HostName: hostName,
		Token:    tokenValue,
	})
	if err != nil {
		return err
	}
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", s.SmtpServerHost, s.SmtpServerPort),
		smtp.PlainAuth("", s.SmtpUser,
			s.SmtpPassword, s.SmtpServerHost),
		s.EmailFrom, []string{email}, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}
