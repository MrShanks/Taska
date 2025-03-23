package utils

import "golang.org/x/crypto/bcrypt"

func Crypt(text *string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*text), 10)
	return string(hash), err
}

func Compare(hashedText, clearText *string) error {
	err := bcrypt.CompareHashAndPassword([]byte(*hashedText), []byte(*clearText))
	if err != nil {
		return err
	}
	return nil
}
