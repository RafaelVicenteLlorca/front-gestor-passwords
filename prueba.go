package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/sha3"
)

var key []byte
var contrasenas []User

type data struct {
	Id uint
	content string
}

type User struct {
	Username string
	Email    string
	Password string
	WebSite  string
	Notes    string
}

func userToString(usuario User) string {
	return "Usuario: " + usuario.Username + " email: " + usuario.Email + " contraseña: " + usuario.Password + " sitioweb: " + usuario.WebSite + " notas: " + usuario.Notes
}

func stringToUser(usuarioString string) User {
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
	user := User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
		WebSite:  strings.TrimSpace(webSite),
		Notes:    strings.TrimSpace(notes),
	}

	return user
}

func generadorHash(email string, pass string) ([]byte, []byte) {
	// Concatenamos el usuario y la contraseña separados por un caracter nulo
	data := []byte(email + "\x00" + pass)

	// Generamos el hash SHA-3 de 512 bits de la concatenación anterior
	hash := sha3.Sum512(data)
	hashemaillargo := sha3.Sum256([]byte(email))
	// Convertimos el hash a una cadena de texto en formato hexadecimal
	// hashStr := hex.EncodeToString(hash[:])

	// Dividimos el hash en dos partes iguales
	key = hash[:len(hash)/2]
	hash2 := hash[len(hash)/2:]
	emailhashed := hashemaillargo[:]
	return hash2, emailhashed
}

func encryptChaCha(texto string) (string, error) {
	// Generar un nonce aleatorio de 24 bytes
	nonce := make([]byte, chacha20.NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Crear un cifrador con la clave y el nonce
	cifrador, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return "", err
	}

	// Cifrar los datos
	cifrado := make([]byte, len(texto))
	cifrador.XORKeyStream(cifrado, []byte(texto))

	// Codificar el nonce y el cifrado como una cadena base64
	resultado := base64.StdEncoding.EncodeToString(append(nonce, cifrado...))

	return resultado, nil
}

func decryptChaCha(ciphertext string) (string, error) {
	// Decodificar la cadena base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Verificar que la longitud de la cadena decodificada sea mayor o igual que la longitud del nonce
	if len(data) < chacha20.NonceSize {
		return "", errors.New("la cadena cifrada es demasiado corta")
	}

	// Extraer el nonce y el cifrado
	nonce := data[:chacha20.NonceSize]
	cifrado := data[chacha20.NonceSize:]

	// Crear un descifrador con la clave y el nonce
	descifrador, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return "", err
	}

	// Descifrar los datos
	plaintext := make([]byte, len(cifrado))
	descifrador.XORKeyStream(plaintext, cifrado)

	return string(plaintext), nil
}

func mainMenu() {
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("---------------------------------------------------------------BIENVENIDO AL GESTOR DE CONTRASEÑAS-------------------------------------------------------------")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Print("\n\n")
	fmt.Println("								Pulse 1 para logearse										")
	fmt.Println("								Pulse 2 para Crear un usuario								")
}

func singinData() ([]byte, []byte) {
	var email string
	var contraseña string
	fmt.Println("Introduzca su email: ")
	fmt.Scan(&email)
	fmt.Println("Introduzca su contraseña: ")
	fmt.Scan(&contraseña)
	login, hashemail := generadorHash(email, contraseña)
	return login, hashemail
}

func singin(loginHash []byte, hashemail []byte) {
	fmt.Println("to do")
}

func singupData() ([]byte, []byte) {
	var email string
	var contraseña string
	var contraseña2 string

	fmt.Println("Introduzca su email: ")
	fmt.Scan(&email)
	for {
		fmt.Println("Introduzca su contraseña, la contraseña debe contener al menos 1 mayuscula y 1 valor numerico: ")
		fmt.Scan(&contraseña)
		if comprobarContrasena(contraseña) {
			fmt.Println("Verificando contraseña, introduzca su contraseña de nuevo: ")
			fmt.Scan(&contraseña2)
			if contraseña == contraseña2 {
				break
			} else {
				fmt.Println("Las contraseñas no coinciden, intentelo de nuevo")
			}
		} else {
			fmt.Println("Error, la contraseña debe contener al menos 1 mayuscula y 1 valor numerico")
		}
	}

	login, hashemail := generadorHash(email, contraseña)
	return login, hashemail
}

func generadorContrasena(tam int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#$%^&*()_+[]{}"
	result := make([]byte, tam)
	for i := range result {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[randomIndex.Int64()]
	}
	return string(result)
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

	var opcion string
	for {
		fmt.Println("Seleccione 1 para una clave aleatoria\n Seleccione 2 para insertar la clave")
		fmt.Scan(&opcion)
		if opcion == "1" {
			tam := 13
			var tamaux string
			var tamaux2 int
			fmt.Println("Seleccione 1 para una clave aleatoria\n Seleccione 2 para insertar la clave")
			fmt.Scan(&tamaux)
			tamaux =
			if   ||tamaux <0 
			password = generadorContrasena(tam)
			break
		} else if opcion == "2" {
			var contrasena string
			var contrasena2 string
			for {
				fmt.Println("Introduzca su contraseña, la contraseña debe contener al menos 1 mayuscula y 1 valor numerico: ")
				fmt.Scan(&contrasena)
				fmt.Println("Verificando contraseña, introduzca la contraseña de nuevo: ")
				fmt.Scan(&contrasena2)
				if comprobarContrasena(contrasena) {
					if contrasena == contrasena2 {
						password = contrasena
						break
					} else {
						fmt.Println("Las contaseñas no coinciden")
					}
				} else {
					fmt.Println("Error, la contraseña debe contener al menos 1 mayuscula y 1 valor numerico")
				}

			}
			break
		} else {
			fmt.Print("La opcion no es correcta, introduzcala de nuevo")
		}
	}
	for {
		fmt.Println("Desea añadir notas (Si/No): ")
		fmt.Scan(&opcion)
		if opcion == "Si" || opcion == "si" {
			fmt.Println("Introduzca las notas : ")
			fmt.Scan(&notes)
			break
		} else if opcion == "No" || opcion == "no" {
			notes = "No hay notas"
			break
		} else {
			fmt.Println("La opcion no es correcta, intentelo de nuevo")
		}
	}
	user := User{
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
			fmt.Println("contaseña " + strconv.Itoa(indice) + ": " + userToString(contrasenas[indice]))
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

	var opcion string
	var nuevo string
	user := contrasenas[posicion]
	fmt.Println("si desaea modificar el nombre de usuario pulse 1\n si desea modificar el email pulse 2\n si desea modificar la contraseña pulse 3\n si desea modificar el sitio web pulse 4\n si desea modificar o añadir notas pulse 5")
	fmt.Scan(&opcion)
	switch opcion {
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
		if opcion == "1" {
			user.Notes = nuevo
		} else if opcion == "2" {
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

func comprobarContrasena(contrasena string) bool {

	if len(contrasena) < 10 {
		return false
	}
	tieneMayuscula := false
	tieneNumero := false
	for _, c := range contrasena {
		if unicode.IsUpper(c) {
			tieneMayuscula = true
		}
		if unicode.IsNumber(c) {
			tieneNumero = true
		}
	}
	valido := tieneMayuscula && tieneNumero

	return valido
}

func pruebas() {
	user := User{
		Username: "vicente",
		Email:    "vicentemail",
		Password: "123",
		WebSite:  "facebook",
		Notes:    "estas notas son de ejemplo",
	}

	user2 := User{
		Username: "vicente2",
		Email:    "vicentemail2",
		Password: "1",
		WebSite:  "facebook2",
		Notes:    "estas notas son de ejemplo2",
	}

	user3 := User{
		Username: "vicente3",
		Email:    "vicentemail3",
		Password: "2",
		WebSite:  "facebook3",
		Notes:    "estas notas son de ejemplo3",
	}
	contrasenas = append(contrasenas, user)
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

	fmt.Println()
	fmt.Println("modificando un usuario")
	modificarContrasena(3)

	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	fmt.Println()
	fmt.Println("borrando un usuario")
	borrarContrasena(2)

	fmt.Println()
	fmt.Println("mostrando usuarios")
	mostrarContrasenas()

	fmt.Println()
	string1 := userToString(user)
	fmt.Println(string1)

	textocifrado, errorsalida := encryptChaCha(string1)
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto cifrado:", textocifrado)

	textodescifrado, errorsalida := decryptChaCha(textocifrado)
	if errorsalida != nil {
		fmt.Println(errorsalida)
	}

	fmt.Println("Texto descifrado:", textodescifrado)

	userX := stringToUser(textodescifrado)
	fmt.Println(userX)
}

func main() {
	var option string
	mainMenu()
	fmt.Scan(&option)
	switch option {
	case "1":
		i := 0
		for {
			loginHash, hashemail := singinData()
			singin(loginHash, hashemail)
			if i == 3 {
				fmt.Println("Numero de intentos superados, intentelo mas tarde")
				break
			}
			i++
			break
		}

	case "2":
		loginHash, hashemail := singupData()
		singin(loginHash, hashemail)
	default:
		fmt.Println("Error al escoger opcion")
	}
	pruebas()
}
