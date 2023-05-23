package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/sha3"
)

func EncodingHashToBase64(hash []byte) string {
	return base64.StdEncoding.EncodeToString(hash)
}

func EncryptChaCha(texto string, key []byte) (string, error) {
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

func DecryptChaCha(ciphertext string, key []byte) (string, error) {
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

func GeneradorHash(email string, pass string) ([]byte, []byte) {
	// Concatenamos el usuario y la contraseña separados por un caracter nulo
	data := []byte(email + "\x00" + pass)

	// Generamos el hash SHA-3 de 512 bits de la concatenación anterior
	hash := sha3.Sum512(data)
	hashemaillargo := sha3.Sum256([]byte(email))
	// Convertimos el hash a una cadena de texto en formato hexadecimal
	// hashStr := hex.EncodeToString(hash[:])

	// Dividimos el hash en dos partes iguales
	hash2 := hash[len(hash)/2:]
	emailhashed := hashemaillargo[:]
	return hash2, emailhashed
}

func RandomPasswordGenerator(size int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#$%^&*()_+[]{}"
	result := make([]byte, size)
	for i := range result {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[randomIndex.Int64()]
	}
	return string(result)
}
