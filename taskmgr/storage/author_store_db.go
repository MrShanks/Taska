package storage

import (
	"context"
	"fmt"

	"github.com/MrShanks/Taska/common/author"
	"github.com/jackc/pgx/v5"
)

type AuthorStore struct {
	Conn *pgx.Conn
}

func (db *AuthorStore) SignUp(newAuthor *author.Author) error {
	query := fmt.Sprintf(
		"INSERT INTO author (firstname, lastname, email, password) values ('%s','%s','%s','%s');",
		newAuthor.Firstname, newAuthor.Lastname, newAuthor.Email, newAuthor.Password)

	_, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("couldn't create author: %v", err)
	}

	return nil
}

func (db *AuthorStore) SignIn(email, password string) error {
	query := fmt.Sprintf("SELECT * FROM author WHERE email = '%s' and password = '%s';", email, password)

	a := author.Author{}

	err := db.Conn.QueryRow(context.Background(), query).Scan(&a.ID, &a.Firstname, &a.Lastname, &a.Email, &a.Password)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("couln't find a email password match")
	}

	if err != nil {
		return fmt.Errorf("error signing in the user: %v", err)
	}

	return nil
}
