package postgresdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://taskauser:secure_password@localhost:5432/taskadb")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func QueryData(conn *pgx.Conn) {
	var id, title, description string
	err := conn.QueryRow(context.Background(), "select * from tasks where id = '4ed7d963-a74c-4231-92dd-f88e2e0b15f9'").Scan(&id, &title, &description)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ID: %s\nTitle: %s\nDescription: %s\n", id, title, description)
	defer conn.Close(context.Background())
}
