package main

import (
	"fmt"
	"os"
	"passwordsAdmin/client"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/services"
	"passwordsAdmin/session"
	loginui "passwordsAdmin/ui/login"
	mainui "passwordsAdmin/ui/main"
	passwordui "passwordsAdmin/ui/password"

	"github.com/TwiN/go-color"
)

const MAX_TRIES_LOGIN = 3

func singInRequest(loginHash string, hashemail string) bool {
	lr := services.LoginRequest{Email: loginHash, Password: hashemail}
	lresp, err := services.UserServiceRequest.Login(lr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	client.HttpClient.SetToken(lresp.Token)
	return true
}

func singUp() (services.RegisterResponse, error) {
	loginHash, hashemail := loginui.SingUpData()
	registerRequest := services.RegisterRequest{Email: loginHash, Password: hashemail}
	return services.UserServiceRequest.Register(registerRequest)
}

func signIn() ([]byte, bool) {
	isSignedIn := false
	var key []byte
	for i := 0; i < MAX_TRIES_LOGIN && !isSignedIn; i++ {
		loginHash, hashemail := loginui.SignInUI()
		if loginHash == nil {
			continue
		}
		key = loginHash
		isSignedIn = singInRequest(utils.EncodingHashToBase64(loginHash), utils.EncodingHashToBase64(hashemail))
	}
	return key, isSignedIn
}

func main() {
	mainui.AppLogo()
	for {
		mainui.MainMenu()
		var option rune
		fmt.Scanf("%c", &option)
		switch string(option) {
		case "1":
			key, isSignedIn := signIn()
			if !isSignedIn {
				fmt.Println(color.Colorize(color.Red, "Numero de intentos superados, inténtelo mas tarde"))
				continue
			}
			fmt.Println(color.Colorize(color.Green, "Logueado"))
			session.SessionObject.SetKey(key)

			passwordui.InitMenu()
		case "2":
			_, err := singUp()
			if err != nil {
				fmt.Println(color.Colorize(color.Red, err.Error()))
				continue
			}
			fmt.Println(color.Colorize(color.Green, "¡Usuario registrado!"))
		case "q":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println(color.Colorize(color.Green, "Error al escoger opcion"))
		}
	}
}
