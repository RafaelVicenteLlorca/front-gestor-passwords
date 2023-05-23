package session

import (
	"passwordsAdmin/client"
	"passwordsAdmin/pkg/user"
)

type Session struct {
	key       []byte
	passwords []user.User
}

var SessionObject = New()

func New() *Session {
	return &Session{key: []byte(""), passwords: []user.User{}}
}

func (s *Session) GetKey() []byte {
	return s.key
}

func (s *Session) SetKey(key []byte) {
	s.key = key
}

func (s *Session) GetPasswords() []user.User {
	return s.passwords
}

func (s *Session) SetPasswords(passwords []user.User) {
	s.passwords = passwords
}

func (s *Session) ClosesSession() error {
	s.SetKey([]byte(""))
	s.SetPasswords([]user.User{})
	client.HttpClient.SetToken("")
	return nil
}
