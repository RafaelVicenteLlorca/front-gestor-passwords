package login

import (
	"fmt"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/session"
	mainui "passwordsAdmin/ui/main"

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

func FormCheckPassword() string {
	passwordInput, err := mainui.RequestPassword("Introduzca su contraseña: ")
	if err != nil {
		return ""
	}
	verifyPassword, err := mainui.RequestPassword("Verificando contraseña, introduzca su contraseña de nuevo: ")
	if err != nil || string(passwordInput) != string(verifyPassword) {
		return ""
	}
	return string(passwordInput)
}

func FormPasswordContinue() string {
	for {
		password := FormCheckPassword()
		if password == "" {
			fmt.Println("Las contraseñas no coinciden, intentelo de nuevo")
			continue
		}
		return password
	}
}

func SingUpData() (string, string) {
	var email string

	fmt.Println("Introduzca su email: ")
	fmt.Scan(&email)
	password := FormPasswordContinue()

	login, hashemail := utils.GeneradorHash(email, string(password))
	session.SessionObject.SetKey(login)
	return utils.EncodingHashToBase64(login), utils.EncodingHashToBase64(hashemail)
}
