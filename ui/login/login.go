package login

import (
	"fmt"
	"prueba/pkg/utils"

	"golang.org/x/term"
)

func SignInUI() ([]byte, []byte) {
	var email string
	fmt.Println("Introduzca su email: ")
	fmt.Scan(&email)
	fmt.Println("Introduzca su contraseña: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println("Error reading password")
		return nil, nil
	}
	return utils.GeneradorHash(email, string(password))
}
