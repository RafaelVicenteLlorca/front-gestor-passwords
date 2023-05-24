package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/services"
)

type User struct {
	Username string
	Email    string
	Password string
	WebSite  string
	Notes    string
}

func (u *User) ToString() string {
	return "Usuario: " + u.Username + ", Email: " + u.Email + ", Contraseña: " + u.Password + ", Sitio Web: " + u.WebSite + ", Notas: " + u.Notes
}

func ConvertPasswordToData(p services.PasswordsResponse, key []byte) (User, error) {
	desencryptedPassword, err := utils.DecryptChaCha(p.Content, key)
	if err != nil {
		return User{}, errors.New("error mientras se actualizaba 1")
	}
	var password User
	err = json.Unmarshal([]byte(desencryptedPassword), &password)
	if err != nil {
		fmt.Println(err.Error())
		return User{}, errors.New("error mientras se actualizaba 2")

	}
	return password, nil
}
