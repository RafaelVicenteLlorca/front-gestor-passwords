package session

import "prueba/pkg/user"

type Session struct {
	key       []byte
	passwords []user.User
}

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
