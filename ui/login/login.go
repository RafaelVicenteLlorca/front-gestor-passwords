package login

import (
	"fmt"
	"os"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/session"
	mainui "passwordsAdmin/ui/main"
	"unicode"

	"github.com/TwiN/go-color"
	"golang.org/x/term"
)

const DEFAULT_MIN_LENGTH = 10

func SignInUI() ([]byte, []byte) {
	var email string
	fmt.Println("Introduzca su email: ")
	fmt.Scanln(&email)
	fmt.Println("Introduzca su contraseña: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
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
			fmt.Println(color.Colorize(color.Red, "Las contraseñas no coinciden, inténtelo de nuevo"))
			continue
		}
		if !CheckSegurityPassword(password) {
			fmt.Println(color.Colorize(color.Red, "Error, la contraseña debe:"))
			fmt.Println(color.Colorize(color.Red, "1. Contener al menos 1 mayúscula"))
			fmt.Println(color.Colorize(color.Red, "2. Valor númerico"))
			fmt.Println(color.Colorize(color.Red, "3. Longitud mínima: "+fmt.Sprintf("%d", DEFAULT_MIN_LENGTH)))
			continue
		}
		return password
	}
}

func SingUpData() (string, string) {
	var email string

	fmt.Println("Introduzca su email: ")
	fmt.Scanln(&email)
	password := FormPasswordContinue()

	login, hashemail := utils.GeneradorHash(email, string(password))
	session.SessionObject.SetKey(login)
	return utils.EncodingHashToBase64(login), utils.EncodingHashToBase64(hashemail)
}
