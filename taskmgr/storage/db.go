package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MrShanks/Taska/common/task"
)

type PostgresDatabase struct {
	Conn *pgx.Conn
}

func (db *PostgresDatabase) Connect(db_url string) error {
	password := os.Getenv("POSTGRES_PWD")
	dburl := fmt.Sprintf(db_url, password)
	var err error
	db.Conn, err = pgx.Connect(context.Background(), dburl)
	if err != nil {
		return err
	}
	return nil
}

func (db *PostgresDatabase) GetOne(id string) (*task.Task, error) {
	query := fmt.Sprintf("select * from tasks where id = '%s';", id)

	t := task.Task{}

	row := db.Conn.QueryRow(context.Background(), query).Scan(&t.ID, &t.Title, &t.Desc)
	if row == pgx.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	fmt.Printf("%v", t)
	return &t, nil
}

func (db *PostgresDatabase) GetTasks() map[uuid.UUID]*task.Task {
	fetchedTasks := make(map[uuid.UUID]*task.Task)
	query := "select * from tasks"

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error to query: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var r task.Task
		err := rows.Scan(&r.ID, &r.Title, &r.Desc)
		if err != nil {
			log.Printf("%v", err)
		}
		fetchedTasks[r.ID] = &task.Task{
			ID:    r.ID,
			Title: r.Title,
			Desc:  r.Desc,
		}
	}
	return fetchedTasks
}

func (db *PostgresDatabase) New(task *task.Task) uuid.UUID {
	task.ID = uuid.New()
	query := fmt.Sprintf("insert into tasks (id, title, description) values ('%s', '%s', '%s');", task.ID, task.Title, task.Desc)

	_, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		log.Printf("%v", err)
		return uuid.Nil
	}
	return task.ID
}

func (db *PostgresDatabase) Update(id, title, desc string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	var query string
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}
	if title != "" {
		query = fmt.Sprintf("update tasks set title = '%v' where id = '%s';", title, id)
	}

	if desc != "" {
		query = fmt.Sprintf("update tasks set description = '%v' where id = '%s';", desc, id)
	}

	update, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error updating task with ID %v with error %v", id, err)
	} else if update.String() == "DELETE 0" {
		return nil, fmt.Errorf("task with ID %v does not exist", UUID)
	}
	return &task.Task{ID: UUID}, nil
}

func (db *PostgresDatabase) Delete(id string) error {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %s", id)
	}
	query := fmt.Sprintf("delete from tasks where id = '%s';", id)

	del, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error deleting task with ID %v with error %v", id, err)
	} else if del.String() == "DELETE 0" {
		return fmt.Errorf("task with ID %v does not exist", UUID)
	}
	log.Printf("Task with ID: %v has been deleted\n", UUID)
	return nil
}
