package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MrShanks/Taska/common/author"
	"github.com/golang-jwt/jwt/v5"
)

var secret string = "secret"

func CreateToken(author author.Author) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": author.Email,
			"iss": "taskmgr",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Second * 30).Unix(),
		})

	singnedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	fmt.Println(singnedToken)
	return singnedToken, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ReadToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("couldn't locate user home\n")
	}

	content, err := os.ReadFile(filepath.Join(home, ".taskcli"))
	if err != nil {
		fmt.Printf("couldn't retrieve token from user home: %v\n", err)
	}

	token := strings.TrimRight(string(content), "\n")

	return token
}

func StoreToken(token string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("couldn't locate user home")
	}

	file, err := os.Create(filepath.Join(home, ".taskcli"))
	if err != nil {
		fmt.Printf("couldn't store login credentials: %v", err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", token))
	if err != nil {
		fmt.Printf("cannot write to %s, error: %v", file.Name(), err)
	}

}
