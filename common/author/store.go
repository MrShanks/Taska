package author

type Store interface {
	SignUp(*Author) error
	SignIn(email, password, token string) error
	GetAuthorID(token string) (string, error)
}
