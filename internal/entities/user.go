package entities

type User struct {
	Id       string
	Email    string
	Password string
	FullName string
	Salt     string
	Role     string
}
