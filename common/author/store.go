package author

type Store interface {
	SignUp(*Author) error
	SignIn(email, password string) error
}
