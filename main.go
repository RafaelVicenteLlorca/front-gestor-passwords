package main

import (
	"fmt"
	"os"
	"passwordsAdmin/client"
	"passwordsAdmin/pkg/user"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/services"
	"passwordsAdmin/session"
	loginui "passwordsAdmin/ui/login"
	mainui "passwordsAdmin/ui/main"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
)

var contrasenas []user.User

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

func generadorContrasena() string {
	// TODO: password auto-generate
	return "to do"
}

func anadirContrasena() {
	var username string
	var email string
	var password string
	var webSite string
	var notes string

	fmt.Println("Introduzca el usuario: ")
	fmt.Scan(&username)

	fmt.Println("Introduzca el email: ")
	fmt.Scan(&email)

	fmt.Println("Introduzca el sitio web: ")
	fmt.Scan(&webSite)

	var opcion []byte = make([]byte, 1)
	for {
		fmt.Println("Seleccione 1 para una clave aleatoria\n Seleccione 2 para insertar la clave")
		os.Stdin.Read(opcion)
		if string(opcion) == "1" {
			password = generadorContrasena()
			break
		} else if string(opcion) == "2" {
			password = loginui.FormPasswordContinue()
			break
		}
		fmt.Print("La opcion no es correcta, introduzcala de nuevo")
	}
	for {
		fmt.Println("Desea añadir notas (s/N): ")
		os.Stdin.Read(opcion)
		if strings.ToLower(string(opcion)) == "s" {
			fmt.Println("Introduzca las notas : ")
			fmt.Scan(&notes)
			break
		} else {
			notes = "No hay notas"
			break
		}
	}
	user1 := user.User{
		Username: username,
		Email:    email,
		Password: password,
		WebSite:  webSite,
		Notes:    notes,
	}

	contrasenas = append(contrasenas, user1)
}

func mostrarContrasenas() {
	indice := 0
	maximo := 0
	var opcion string
	for {
		maximo += 5
		if maximo > (len(contrasenas) - 1) {
			maximo = len(contrasenas) - 1
		}
		for {
			fmt.Println("contaseña " + strconv.Itoa(indice) + ": " + contrasenas[indice].ToString())
			if maximo == indice {
				break
			}
			indice++
		}
		if indice < (len(contrasenas) - 1) {
			fmt.Println("si desea ver las siguiente 5 introduzca 1 si desea ver los 5 anteriores introduzca 2 si no desea ver mas contraseñas introduzca 3")
			fmt.Scan(&opcion)
			if opcion == "2" {
				indice -= 5
			} else if opcion == "3" {
				break
			}
		} else {
			break
		}
	}
}

func modificarContrasena(posicion int) {

	var opcion []byte = make([]byte, 1)
	var nuevo string
	user1 := contrasenas[posicion]
	fmt.Println("si desaea modificar el nombre de usuario pulse 1\n si desea modificar el email pulse 2\n si desea modificar la contraseña pulse 3\n si desea modificar el sitio web pulse 4\n si desea modificar o añadir notas pulse 5")
	os.Stdin.Read(opcion)
	switch string(opcion) {
	case "1":
		fmt.Println("Introduzca el nuevo nombre: ")
		fmt.Scan(&nuevo)
		user1.Username = nuevo
		contrasenas[posicion] = user1
	case "2":
		fmt.Println("Introduzca el nuevo email: ")
		fmt.Scan(&nuevo)
		user1.Email = nuevo
		contrasenas[posicion] = user1
	case "3":
		fmt.Println("Introduzca la nueva contraseña: ")
		fmt.Scan(&nuevo)
		user1.Password = nuevo
		contrasenas[posicion] = user1
	case "4":
		fmt.Println("Introduzca el nuevo sitio web: ")
		fmt.Scan(&nuevo)
		user1.WebSite = nuevo
		contrasenas[posicion] = user1
	case "5":
		fmt.Println("desea añadir otra nota o modificar la actual (Añadir 1, Modificar 2): ")
		fmt.Scan(&opcion)
		fmt.Println("Introduzca la nueva nota: ")
		fmt.Scan(&nuevo)
		if string(opcion) == "1" {
			user1.Notes = nuevo
		} else if string(opcion) == "2" {
			user1.Notes += ", " + nuevo
		}
		contrasenas[posicion] = user1
	default:
		fmt.Println("el valor no es correcto")
	}
}

func borrarContrasena(posicion int) {
	contrasenas = append(contrasenas[:posicion], contrasenas[posicion+1:]...)
}

func passwordsAdmin() {
	user1 := user.User{
		Username: "vicente",
		Email:    "vicentemail",
		Password: "123",
		WebSite:  "facebook",
		Notes:    "estas notas son de ejemplo",
	}

	user2 := user.User{
		Username: "vicente2",
		Email:    "vicentemail2",
		Password: "1",
		WebSite:  "facebook2",
		Notes:    "estas notas son de ejemplo2",
	}

	user3 := user.User{
		Username: "vicente3",
		Email:    "vicentemail3",
		Password: "2",
		WebSite:  "facebook3",
		Notes:    "estas notas son de ejemplo3",
	}
	contrasenas = append(contrasenas, user1)
	contrasenas = append(contrasenas, user2)
	contrasenas = append(contrasenas, user3)
	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	fmt.Println()
	fmt.Println("añadirendo un usuario")
	anadirContrasena()

	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	/* fmt.Println()
	fmt.Println("modificando un usuario")
	modificarContrasena(3) */

	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	/* fmt.Println()
	fmt.Println("borrando un usuario")
	borrarContrasena(2) */

	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	fmt.Println()
	string1 := user1.ToString()
	fmt.Println(string1)

	textocifrado, errorsalida := utils.EncryptChaCha(string1, session.SessionObject.GetKey())
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto cifrado:", textocifrado)

	textodescifrado, errorsalida := utils.DecryptChaCha(textocifrado, session.SessionObject.GetKey())
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto descifrado:", textodescifrado)

	userX := user.NewByString(textodescifrado)
	fmt.Println(userX)
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
		var option []byte = make([]byte, 1)
		os.Stdin.Read(option)
		switch string(option) {
		case "1":
			key, isSignedIn := signIn()
			if !isSignedIn {
				fmt.Println(color.Colorize(color.Red, "Numero de intentos superados, inténtelo mas tarde"))
				continue
			}
			fmt.Println(color.Colorize(color.Green, "Logueado"))
			session.SessionObject.SetKey(key)

			// TODO: go to initial page
		case "2":
			_, err := singUp()
			if err != nil {
				fmt.Println(color.Colorize(color.Red, err.Error()))
				continue
			}
			fmt.Println("usuario registrado!")
		case "q":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("Error al escoger opcion")
		}
	}
}
