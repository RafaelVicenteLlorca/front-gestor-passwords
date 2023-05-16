package main

import (
	"fmt"
	"os"
	"prueba/client"
	"prueba/pkg/user"
	"prueba/pkg/utils"
	"prueba/services"
	"prueba/session"
	loginui "prueba/ui/login"
	mainui "prueba/ui/main"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var contrasenas []user.User
var sessionObject *session.Session

var (
	userService *services.UserService
)

const MAX_TRIES_LOGIN = 3

func stringToUser(usuarioString string) user.User {
	// Dividir la cadena en tres partes usando el separador " "

	partes := strings.Split(usuarioString, " ")

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
	user := user.User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
		WebSite:  strings.TrimSpace(webSite),
		Notes:    strings.TrimSpace(notes),
	}

	return user
}

func singin(loginHash string, hashemail string) bool {
	lr := services.LoginRequest{Email: loginHash, Password: hashemail}
	lresp, err := userService.Login(lr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	userService.HttpClient.SetToken(lresp.Token)
	return true
}

func formPasswordContinue() string {
	for {
		password := formPassword()
		if password == "" {
			fmt.Println("Las contraseñas no coinciden, intentelo de nuevo")
			continue
		}
		return password
	}
}

func formPassword() string {
	fmt.Println("Introduzca su contraseña: ")
	passwordInput, err := term.ReadPassword(0)
	if err != nil {
		return ""
	}
	fmt.Println("Verificando contraseña, introduzca su contraseña de nuevo: ")
	verifyPassword, err := term.ReadPassword(0)
	if err != nil {
		return ""
	}
	if string(passwordInput) != string(verifyPassword) {
		return ""
	}
	return string(passwordInput)
}

func singupData() (string, string) {
	var email string

	fmt.Println("Introduzca su email: ")
	fmt.Scan(&email)
	password := formPasswordContinue()

	login, hashemail := utils.GeneradorHash(email, string(password))
	sessionObject.SetKey(login)
	return utils.EncodingHashToBase64(login), utils.EncodingHashToBase64(hashemail)
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
			password = formPasswordContinue()
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
	user := user.User{
		Username: username,
		Email:    email,
		Password: password,
		WebSite:  webSite,
		Notes:    notes,
	}

	contrasenas = append(contrasenas, user)
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
			fmt.Println("contaseña " + strconv.Itoa(indice) + ": " + contrasenas[indice].UserToString())
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
	user := contrasenas[posicion]
	fmt.Println("si desaea modificar el nombre de usuario pulse 1\n si desea modificar el email pulse 2\n si desea modificar la contraseña pulse 3\n si desea modificar el sitio web pulse 4\n si desea modificar o añadir notas pulse 5")
	os.Stdin.Read(opcion)
	switch string(opcion) {
	case "1":
		fmt.Println("Introduzca el nuevo nombre: ")
		fmt.Scan(&nuevo)
		user.Username = nuevo
		contrasenas[posicion] = user
	case "2":
		fmt.Println("Introduzca el nuevo email: ")
		fmt.Scan(&nuevo)
		user.Email = nuevo
		contrasenas[posicion] = user
	case "3":
		fmt.Println("Introduzca la nueva contraseña: ")
		fmt.Scan(&nuevo)
		user.Password = nuevo
		contrasenas[posicion] = user
	case "4":
		fmt.Println("Introduzca el nuevo sitio web: ")
		fmt.Scan(&nuevo)
		user.WebSite = nuevo
		contrasenas[posicion] = user
	case "5":
		fmt.Println("desea añadir otra nota o modificar la actual (Añadir 1, Modificar 2): ")
		fmt.Scan(&opcion)
		fmt.Println("Introduzca la nueva nota: ")
		fmt.Scan(&nuevo)
		if string(opcion) == "1" {
			user.Notes = nuevo
		} else if string(opcion) == "2" {
			user.Notes += ", " + nuevo
		}
		contrasenas[posicion] = user
	default:
		fmt.Println("el valor no es correcto")
	}
}

func borrarContrasena(posicion int) {
	contrasenas = append(contrasenas[:posicion], contrasenas[posicion+1:]...)
}

func pruebas() {
	user := user.User{
		Username: "vicente",
		Email:    "vicentemail",
		Password: "123",
		WebSite:  "facebook",
		Notes:    "estas notas son de ejemplo",
	}

	/* user2 := user.User{
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
	} */
	contrasenas = append(contrasenas, user)
	/* contrasenas = append(contrasenas, user2)
	contrasenas = append(contrasenas, user3) */
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
	string1 := user.UserToString()
	fmt.Println(string1)

	textocifrado, errorsalida := utils.EncryptChaCha(string1, sessionObject.GetKey())
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto cifrado:", textocifrado)

	textodescifrado, errorsalida := utils.DecryptChaCha(textocifrado, sessionObject.GetKey())
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto descifrado:", textodescifrado)

	userX := stringToUser(textodescifrado)
	fmt.Println(userX)
}

func initServices() {
	httpClient := client.New()
	userService = &services.UserService{HttpClient: httpClient}
}

func main() {
	sessionObject = session.New()
	initServices()
	mainui.AppLogo()
	for {
		mainui.MainMenu()
		var option []byte = make([]byte, 1)
		os.Stdin.Read(option)
		switch string(option) {
		case "1":
			isSignedIn := false
			var key []byte
			for i := 0; i < MAX_TRIES_LOGIN && !isSignedIn; i++ {
				loginHash, hashemail := loginui.SignInUI()
				if loginHash == nil {
					continue
				}
				key = loginHash
				isSignedIn = singin(utils.EncodingHashToBase64(loginHash), utils.EncodingHashToBase64(hashemail))
			}

			if !isSignedIn {
				fmt.Println("Numero de intentos superados, inténtelo mas tarde")
			} else {
				fmt.Println("Logueado")
				sessionObject.SetKey(key)

				// TODO: go to initial page
			}

		case "2":
			loginHash, hashemail := singupData()
			fmt.Println(loginHash)
			fmt.Println(hashemail)
			// TODO: userService register user

			// TODO: go to main menu
		case "q":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("Error al escoger opcion")
		}
	}
	//pruebas()
}
