package login

import (
	"fmt"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/session"
	mainui "passwordsAdmin/ui/main"
	"unicode"

	"golang.org/x/term"
)

const DEFAULT_MIN_LENGTH = 10

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
	passwordInput, err := mainui.RequestPassword("Introduzca su contraseña, la contraseña debe contener al menos 1 mayúscula y 1 valor númerico: ")
	if err != nil {
		return ""
	}
	verifyPassword, err := mainui.RequestPassword("Verificando contraseña, introduzca su contraseña de nuevo: ")
	if err != nil || string(passwordInput) != string(verifyPassword) {
		return ""
	}
	return string(passwordInput)
}

func CheckSegurityPassword(password string) bool {
	if len(password) < DEFAULT_MIN_LENGTH {
		return false
	}
	containsUppercase := false
	containsNumber := false
	for _, c := range password {
		containsUppercase = containsUppercase || unicode.IsUpper(c)
		containsNumber = containsNumber || unicode.IsNumber(c)
		if containsUppercase && containsNumber {
			break
		}
	}
	return containsUppercase && containsNumber
}

func FormPasswordContinue() string {
	for {
		password := FormCheckPassword()
		if password == "" {
			fmt.Println("Las contraseñas no coinciden, intentelo de nuevo")
			continue
		}
		if CheckSegurityPassword(password) {
			fmt.Println("Error, la contraseña debe contener al menos 1 mayúscula y 1 valor númerico")

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
