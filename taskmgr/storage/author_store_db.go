package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/MrShanks/Taska/common/author"
	"github.com/MrShanks/Taska/utils"
	"github.com/jackc/pgx/v5"
)

type AuthorStore struct {
	Conn *pgx.Conn
}

func (db *AuthorStore) GetAuthorID(token string) (string, error) {
	query := fmt.Sprintf("SELECT id FROM author WHERE token = '%s';", token)

	var authorID string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := db.Conn.QueryRow(ctx, query).Scan(&authorID)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("no author found with token: %s. Error: %v", token, err)
	}

	return authorID, nil
}

func (db *AuthorStore) SignUp(newAuthor *author.Author) error {
	var err error

	newAuthor.Password, err = utils.Crypt(&newAuthor.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	query := fmt.Sprintf(
		"INSERT INTO author (firstname, lastname, email, password) values ('%s','%s','%s','%s');",
		newAuthor.Firstname, newAuthor.Lastname, newAuthor.Email, newAuthor.Password)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = db.Conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("couldn't create author: %v", err)
	}

	return nil
}

func (db *AuthorStore) SignIn(email, password, token string) error {
	query := fmt.Sprintf("SELECT * FROM author WHERE email = '%s'", email)

	a := author.Author{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := db.Conn.QueryRow(ctx, query).Scan(&a.ID, &a.Firstname, &a.Lastname, &a.Email, &a.Password, &a.Token)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("couln't find an email match")
	}

	if utils.Compare(&a.Password, &password) != nil {
		return fmt.Errorf("incorrect password: %v", err)
	}

	if err != nil {
		return fmt.Errorf("error signing in the user: %v", err)
	}

	query = fmt.Sprintf("UPDATE author SET token = '%v' WHERE email = '%s'", token, email)

	_, err = db.Conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("couldn't save author token: %v", err)
	}

	return nil
}
