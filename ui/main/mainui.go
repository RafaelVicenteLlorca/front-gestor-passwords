package mainui

import (
	"fmt"

	"golang.org/x/term"
)

func AppLogo() {
	fmt.Println()
	fmt.Println("        #####              ########      ")
	fmt.Println("     ###########         #####  #####    ")
	fmt.Println("  ###############      ###         ####  ")
	fmt.Println(" ####         ###     ####          #### ")
	fmt.Println(" ###           ##     ####           ### ")
	fmt.Println("###             #     ####           ##  ")
	fmt.Println("###                   ####           #   ")
	fmt.Println("###                   ####          #    ")
	fmt.Println("###                   ####        ##     ")
	fmt.Println("###          ####     ####      ###      ")
	fmt.Println("###           ###     ####  ####         ")
	fmt.Println("###             #     #######            ")
	fmt.Println("###             #     ####               ")
	fmt.Println("####          ###     ####               ")
	fmt.Println(" ####       ####      ####               ")
	fmt.Println(" ##############       ####               ")
	fmt.Println("  ############        ####               ")
	fmt.Println()
	fmt.Println("GRAND PASSWORD")
}

func MainMenu() {
	fmt.Println()
	fmt.Println("--------------------------------")
	fmt.Println("1 Login")
	fmt.Println("2 Crear usuario")
	fmt.Println("(q para cerrar)")
}

func RequestPassword(text string) ([]byte, error) {
	fmt.Print(text)
	password, err := term.ReadPassword(0)
	fmt.Println()
	return password, err
}
