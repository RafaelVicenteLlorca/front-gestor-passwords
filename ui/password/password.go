package password

import (
	"encoding/json"
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
	content, err := json.Marshal(password)
	if err != nil {
		return err
	}
	encriptedContent, err := utils.EncryptChaCha(string(content), session.SessionObject.GetKey())
	if err != nil {
		return err
	}
	body := services.PasswordsCreateRequest{Content: encriptedContent}
	newPassword, err := services.PasswordsServiceRequest.Create(body)
	if err != nil {
		return err
	}
	session.SessionObject.SetPasswords(append(session.SessionObject.GetPasswords(), *newPassword))
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
		fmt.Scanln(&newPassword.Notes)
	} else {
		newPassword.Notes = "No hay notas"
	}

	return newPassword
}

func GetAllPasswords() error {
	fmt.Println(color.Colorize(color.Blue, "Cargando contraseñas..."))
	fmt.Println()
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
			var p user.User
			err := json.Unmarshal([]byte(desencryptedPassword), &p)
			if err != nil {
				fmt.Println(color.Colorize(color.Red, "Contraseña mal serializada"))
			}
			fmt.Println(p.ToString())
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
	var newPassword = user.User{}

	fmt.Print("Introduzca el usuario: ")
	fmt.Printf("[%s]", password.Username)
	fmt.Scanln(&newPassword.Username)
	fmt.Println()
	if newPassword.Username == "" {
		newPassword.Username = password.Username
	}

	fmt.Print("Introduzca el email: ")
	fmt.Printf("[%s]", password.Email)
	fmt.Scanln(&newPassword.Email)
	fmt.Println()
	if newPassword.Email == "" {
		newPassword.Email = password.Email
	}

	fmt.Print("Introduzca el sitio web: ")
	fmt.Printf("[%s]", password.WebSite)
	fmt.Scanln(&newPassword.WebSite)
	fmt.Println()
	if newPassword.WebSite == "" {
		newPassword.WebSite = password.WebSite
	}

	var opcion string
	for {
		fmt.Println("Seleccione 1 para una clave aleatoria")
		fmt.Println("Seleccione 2 para insertar la clave")
		fmt.Println("Otro para mantener la contraseña actual")
		fmt.Scanln(&opcion)
		if opcion == "1" {
			var tamaux int
			fmt.Printf("Introduzca longitud de la contraseña [%d]:", default_password_size)
			fmt.Scanln(&tamaux)
			if tamaux == 0 {
				tamaux = default_password_size
			}
			newPassword.Password = utils.RandomPasswordGenerator(tamaux)
			break
		}
		if opcion == "2" {
			password, err := mainui.RequestPassword("Contraseña: ")
			if err != nil {
				fmt.Print(color.Colorize(color.Red, "ERROR con la contraseña, repite"))
				continue
			}
			newPassword.Password = string(password)
			break
		}
		newPassword.Password = password.Password
		break
	}
	fmt.Println("Desea añadir notas (s/N): ")
	fmt.Scanln(&opcion)

	if strings.ToLower(string(opcion)) == "s" {
		fmt.Println("Introduzca las notas: ")
		fmt.Scanln(&newPassword.Notes)
	} else {
		newPassword.Notes = password.Notes
	}

	return &newPassword, nil
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
	password, err := user.ConvertPasswordToData(*passwordSelected, session.SessionObject.GetKey())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	newPassword, err := updatePassword(password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	content, err := json.Marshal(newPassword)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	encriptedContent, err := utils.EncryptChaCha(string(content), session.SessionObject.GetKey())
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(color.Colorize(color.Red, "Error al actualizar la contraseña"))
		return
	}
	body := services.PasswordsUpdateRequest{Content: encriptedContent}
	err = services.PasswordsServiceRequest.Update(passwordSelected.ID, body)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(color.Colorize(color.Red, "Error al actualizar la contraseña"))
		return
	}

	var passwords = session.SessionObject.GetPasswords()
	passwords[option-1] = services.PasswordsResponse{ID: passwordSelected.ID, Content: encriptedContent}
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
