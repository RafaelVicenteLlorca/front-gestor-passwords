package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/sha3"
)

var key []byte

func generadorhas(nombre string, pass string) []byte {
	// Concatenamos el usuario y la contraseña separados por un caracter nulo
	data := []byte(nombre + "\x00" + pass)

	// Generamos el hash SHA-3 de 512 bits de la concatenación anterior
	hash := sha3.Sum512(data)

	// Convertimos el hash a una cadena de texto en formato hexadecimal
	// hashStr := hex.EncodeToString(hash[:])

	// Dividimos el hash en dos partes iguales
	key = hash[:len(hash)/2]
	hash2 := hash[len(hash)/2:]

	return hash2
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

func main() {

	var nombre string
	var contraseña string
	fmt.Println("nombre: ")
	fmt.Scan(&nombre)
	fmt.Println("pass: ")
	fmt.Scan(&contraseña)

	fmt.Println("Mi nombre es", nombre, "y la contraseña es", contraseña, ".")

	hashStr := generadorhas(nombre, contraseña)
	fmt.Println("Hash generado:", hex.EncodeToString(hashStr[:]))

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
}
