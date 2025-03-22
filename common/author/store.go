package author

type Store interface {
	SignUp(*Author) error
	SignIn(email, password string) error
	Crypt(text *string) (string, error)
	Compare(hashedText, clearText *string) error
}
