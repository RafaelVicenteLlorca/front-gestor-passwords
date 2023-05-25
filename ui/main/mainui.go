package mainui

import (
	"fmt"
	"os"

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
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return password, err
}
