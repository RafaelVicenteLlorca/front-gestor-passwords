package mainui

import (
	"fmt"

	"golang.org/x/term"
)

func AppLogo() {
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("---------------------------------------------------------------BIENVENIDO AL GESTOR DE CONTRASEÑAS-------------------------------------------------------------")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------")
}

func MainMenu() {
	fmt.Print("\n\n")
	fmt.Println("								Pulse 1 para logearse										")
	fmt.Println("								Pulse 2 para Crear un usuario								")
	fmt.Println("(Pulse q para cerrar)")
}

func RequestPassword(text string) ([]byte, error) {
	fmt.Print(text)
	password, err := term.ReadPassword(0)
	fmt.Println()
	return password, err
}
