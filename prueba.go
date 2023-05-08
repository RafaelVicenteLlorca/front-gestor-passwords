package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/sha3"
)

var key []byte

type User struct {
	Username string
	Email    string
	Password string
}

func UserToString(usuario User) string {
	return "Usuario: " + usuario.Username + " email: " + usuario.Email + " contraseña: " + usuario.Password
}

func StringToUser(linea string) User {
	// Dividir la cadena en tres partes usando el separador " "
	partes := strings.Split(linea, " ")

	// Extraer los valores de usuario, email y contraseña
	username := strings.TrimPrefix(partes[0], "Usuario:")
	email := strings.TrimPrefix(partes[1], "email:")
	password := strings.TrimPrefix(partes[2], "contraseña:")

	// Crear una nueva instancia de la estructura User con los valores extraídos
	user := &User{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
	}

	return user
}

func generadorhas(email string, pass string) ([]byte, []byte) {
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

func EncryptChaCha(texto string) (string, error) {
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

func DecryptChaCha(ciphertext string) (string, error) {
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

func mainmenu() {
	fmt.Println("------------------------------------------------------------------------------------------/n")
	fmt.Println("-------------------------------BIENVENIDO AL GESTOR DE CONTRASEÑAS------------------------/n")
	fmt.Println("------------------------------------------------------------------------------------------/n")
	fmt.Println("/n /n /n")
	fmt.Println("								Pulse 1 para logearse										/n")
	fmt.Println("								Pulse 2 para Crear un usuario								/n")
}

func singindata() ([]byte, []byte) {
	var email string
	var contraseña string
	fmt.Println("introduzca su email: ")
	fmt.Scan(&email)
	fmt.Println("introduzca su contraseña: ")
	fmt.Scan(&contraseña)
	login, hashemail := generadorhas(email, contraseña)
	return login, hashemail
}

func singin([]byte, []byte) {
	fmt.Println("to do")
}

func singupdata() ([]byte, []byte) {
	var email string
	var contraseña string
	var contraseña2 string

	fmt.Println("introduzca su email: ")
	fmt.Scan(&email)
	for {
		fmt.Println("introduzca su contraseña: ")
		fmt.Scan(&contraseña)
		fmt.Println("verificando contraseña, introduzca su contraseña de nuevo: ")
		fmt.Scan(&contraseña2)
		if contraseña == contraseña2 {
			break
		}
	}

	login, hashemail := generadorhas(email, contraseña)
	return login, hashemail
}

func main() {
	var option string
	mainmenu()
	fmt.Scan(&option)
	switch option {
	case "1":
		i := 0
		for {
			loginHash, hashemail := singindata()
			singin(loginHash, hashemail)
			if i == 3 {
				fmt.Println("numero de intentos superados, intentelo mas tarde")
			}
		}

	case "2":
		loginHash, hashemail := singupdata()
		singin(loginHash, hashemail)
	default:
		fmt.Println("error al escoger opcion")
	}
	/*
		fmt.Println("Hash generado:", hex.EncodeToString(loginHash[:]))

		textocifrado, errorsalida := EncryptChaCha("hola mundo")
		if errorsalida != nil {
			fmt.Println(errorsalida)
		}

		fmt.Println("texto cifrado:", textocifrado)

		textodescifrado, errorsalida := DecryptChaCha(textocifrado)
		if errorsalida != nil {
			fmt.Println(errorsalida)
		}

		fmt.Println("texto descifrado:", textodescifrado)
	*/
}
