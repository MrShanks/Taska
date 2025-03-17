package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
