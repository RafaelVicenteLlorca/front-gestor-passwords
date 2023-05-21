package user

import (
	"strings"
)

type User struct {
	Username string
	Email    string
	Password string
	WebSite  string
	Notes    string
}

func (u *User) ToString() string {
	return "Usuario: " + u.Username + " email: " + u.Email + " contraseña: " + u.Password + " sitioweb: " + u.WebSite + " notas: " + u.Notes
}

func NewByString(s string) *User {
	// Dividir la cadena en tres partes usando el separador " "

	partes := strings.Split(s, " ")

	// Extraer los valores de usuario, email y contraseña
	username := strings.TrimPrefix(partes[1], "Usuario:")
	email := strings.TrimPrefix(partes[3], "email:")
	password := strings.TrimPrefix(partes[5], "contraseña:")
	webSite := strings.TrimPrefix(partes[7], "sitioweb:")

	notesIndex := -1

	for i, v := range partes {
		if v == "notas:" {
			notesIndex = i
			break
		}
	}

	var notes string

	if notesIndex != -1 {
		notes = strings.Join(partes[notesIndex+1:], " ")
	}
	// Crear una nueva instancia de la estructura User con los valores extraídos
	return &User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
		WebSite:  strings.TrimSpace(webSite),
		Notes:    strings.TrimSpace(notes),
	}
}
