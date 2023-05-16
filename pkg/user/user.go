package user

type User struct {
	Username string
	Email    string
	Password string
	WebSite  string
	Notes    string
}

func (u *User) ToString() string {
	return "Usuario: " + u.Username + " email: " + u.Email + " contraseña: " + u.Password + " sitioweb: " + u.WebSite + " notas: " + u.Notes
}
