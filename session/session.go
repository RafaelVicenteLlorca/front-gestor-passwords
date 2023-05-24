package session

import (
	"errors"
	"passwordsAdmin/client"
	"passwordsAdmin/services"
)

type Session struct {
	key       []byte
	passwords []services.PasswordsResponse
}

var SessionObject = New()

func New() *Session {
	return &Session{key: []byte(""), passwords: []services.PasswordsResponse{}}
}

func (s *Session) GetKey() []byte {
	return s.key
}

func (s *Session) SetKey(key []byte) {
	s.key = key
}

func (s *Session) GetPasswords() []services.PasswordsResponse {
	return s.passwords
}

func (s *Session) SetPasswords(passwords []services.PasswordsResponse) {
	s.passwords = passwords
}

func (s *Session) GetPasswordByPosition(i int) (*services.PasswordsResponse, error) {
	if i < 0 || i >= len(s.passwords) {
		return nil, errors.New("IndexedError")
	}

	return &s.passwords[i], nil
}

func (s *Session) DeletePassword(i int) error {
	if i < 0 || i >= len(s.passwords) {
		return errors.New("IndexedError")
	}
	s.passwords = append(s.passwords[:i], s.passwords[i+1:]...)
	return nil
}

func (s *Session) ClosesSession() error {
	s.SetKey([]byte(""))
	s.SetPasswords([]services.PasswordsResponse{})
	client.HttpClient.SetToken("")
	return nil
}
