package password

import (
	"encoding/json"
	"errors"
	"fmt"
	"passwordsAdmin/pkg/user"
	"passwordsAdmin/pkg/utils"
	"passwordsAdmin/services"
	"passwordsAdmin/session"
	mainui "passwordsAdmin/ui/main"
	"strings"

	"github.com/TwiN/go-color"
)

const default_password_size = 13

func AddPassword() error {
	password := AddPasswordUI()
	content := fmt.Sprintf("%+v", password)
	encriptedContent, err := utils.EncryptChaCha(content, session.SessionObject.GetKey())
	if err != nil {
		return err
	}
	body := services.PasswordsCreateRequest{Content: encriptedContent}
	services.PasswordsServiceRequest.Create(body)
	return nil
}

func AddPasswordUI() user.User {
	var newPassword = user.User{}

	fmt.Println("Introduzca el usuario: ")
	fmt.Scan(&newPassword.Username)

	fmt.Println("Introduzca el email: ")
	fmt.Scan(&newPassword.Email)

	fmt.Println("Introduzca el sitio web: ")
	fmt.Scan(&newPassword.WebSite)

	var opcion string
	for {
		fmt.Println("Seleccione 1 para una clave aleatoria")
		fmt.Println("Seleccione 2 para insertar la clave")
		fmt.Scan(&opcion)
		if string(opcion) == "1" {
			var tamaux int
			fmt.Printf("Introduzca longitud de la contraseña [%d]:", default_password_size)
			fmt.Scanln(&tamaux)
			if tamaux == 0 {
				tamaux = default_password_size
			}
			newPassword.Password = utils.RandomPasswordGenerator(tamaux)
			break
		} else if string(opcion) == "2" {
			password, err := mainui.RequestPassword("Contraseña: ")
			if err != nil {
				fmt.Print(color.Colorize(color.Red, "ERROR con la contraseña, repite"))
				continue
			}
			newPassword.Password = string(password)
			break
		}
		fmt.Println("La opcion no es correcta, introduzcala de nuevo")
	}
	fmt.Println("Desea añadir notas (s/N): ")
	fmt.Scanln(&opcion)

	if strings.ToLower(string(opcion)) == "s" {
		fmt.Println("Introduzca las notas : ")
		fmt.Scan(&newPassword.Notes)
	} else {
		newPassword.Notes = "No hay notas"
	}

	return newPassword
}

func GetAllPasswords() error {
	fmt.Println(color.Colorize(color.Blue, "Cargando contraseñas..."))
	passwords, err := services.PasswordsServiceRequest.GetAll()
	if err != nil {
		return err
	}
	session.SessionObject.SetPasswords(passwords)
	return nil
}

func showPasswords() error {
	for i, password := range session.SessionObject.GetPasswords() {
		fmt.Printf("%d: ", i+1)
		desencryptedPassword, err := utils.DecryptChaCha(password.Content, session.SessionObject.GetKey())
		if err != nil {
			fmt.Println(color.Colorize(color.Red, err.Error()))
		} else {
			fmt.Println(desencryptedPassword)
		}
	}
	return nil
}

func deletePasswordsUI(password services.PasswordsResponse) bool {
	var option string
	fmt.Println("¿Estás seguro de eliminar la contraseña? (s/N)")
	fmt.Scanln(&option)
	if strings.ToLower(option) != "s" {
		return false
	}
	err := services.PasswordsServiceRequest.Delete(password.ID)
	if err != nil {
		fmt.Println(color.Colorize(color.Red, "Error al borrar la contraseña"))
		return false
	}
	return true
}

func updatePassword(password user.User) (*user.User, error) {
	// TODO: add update form and request to update by id
	/* var opcion []byte = make([]byte, 1)
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
	} */
	return &user.User{}, nil
}

func convertPasswordToData(p services.PasswordsResponse) (user.User, error) {
	desencryptedPassword, err := utils.DecryptChaCha(p.Content, session.SessionObject.GetKey())
	if err != nil {
		return user.User{}, errors.New("error mientras se actualizaba")
	}
	var password user.User
	err = json.Unmarshal([]byte(desencryptedPassword), &password)
	if err != nil {
		return user.User{}, errors.New("error mientras se actualizaba")

	}
	return password, nil
}

func InitMenu() {
	GetAllPasswords()
	for {
		err := showPasswords()
		if err != nil {
			break
		}
		option := PasswordActions()

		if option == "c" {
			fmt.Println(color.Colorize(color.Blue, "Cerrando sesión..."))
			session.SessionObject.ClosesSession()
			break
		}
		if option == "1" {
			optionAdd()
			continue
		}
		if option == "2" {
			optionGetAll()
			continue
		}

		if option == "3" {
			optionUpdate()
			continue
		}
		if option == "4" {
			optionDelete()
			continue
		}

	}
}

func optionAdd() {
	err := AddPassword()
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
	}
}

func optionGetAll() {
	err := GetAllPasswords()
	if err != nil {
		fmt.Println(color.Colorize(color.Red, err.Error()))
	}
}

func optionUpdate() {
	var passwordSelected *services.PasswordsResponse
	var option int
	for {
		option = requestByPassword()
		password, err := session.SessionObject.GetPasswordByPosition(option - 1)
		if err != nil {
			fmt.Println(color.Colorize(color.Red, "Contraseña no encontrada"))
			continue
		}
		passwordSelected = password
		break
	}
	// transform password to user
	password, err := convertPasswordToData(*passwordSelected)
	if err != nil {
		return
	}
	newPassword, err := updatePassword(password)
	if err != nil {
		return
	}
	fmt.Println(newPassword)
	// ! call to api

	// * result: if is true, updated in array, false, exit

}

func optionDelete() {
	var passwordSelected *services.PasswordsResponse
	var option int
	for {
		option = requestByPassword()
		password, err := session.SessionObject.GetPasswordByPosition(option - 1)
		if err != nil {
			fmt.Println(color.Colorize(color.Red, "Contraseña no encontrada"))
			continue
		}
		passwordSelected = password
		break
	}
	result := deletePasswordsUI(*passwordSelected)
	if result {
		session.SessionObject.DeletePassword(option - 1)
		fmt.Println(color.Colorize(color.Green, "¡Eliminada correctamente!"))
	}
}

func requestByPassword() int {
	var option int
	fmt.Print("Introduce el número de la contraseña a seleccionar: ")
	fmt.Scanf("%d", &option)
	return option
}

func PasswordActions() string {
	var option string
	fmt.Println()
	fmt.Println()
	fmt.Println("-------------------------")
	fmt.Println("1 Añadir contraseñas")
	fmt.Println("2 Ver todas las contraseñas")
	fmt.Println("3 Editar contraseña")
	fmt.Println("4 Borrar contraseña")
	fmt.Println("(c para cerrar sesión)")
	fmt.Scan(&option)
	return option
}
