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
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	return nil
}

func (db *PostgresDatabase) GetOne(id string) (*task.Task, error) {
	query := fmt.Sprintf("select * from task where id = '%s';", id)

	t := task.Task{}

	row := db.Conn.QueryRow(context.Background(), query).Scan(&t.ID, &t.Title, &t.Desc)
	if row == pgx.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}

	return &t, nil
}

func (db *PostgresDatabase) GetTasks() []*task.Task {
	var fetchedTasks []*task.Task
	query := "select * from task"

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error to query: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		r := &task.Task{}

		err := rows.Scan(&r.ID, &r.Title, &r.Desc)
		if err != nil {
			log.Printf("%v", err)
		}

		fetchedTasks = append(fetchedTasks, r)
	}
	return fetchedTasks
}

func (db *PostgresDatabase) New(task *task.Task) uuid.UUID {
	query := fmt.Sprintf("insert into task (title, description) values ('%s', '%s');", task.Title, task.Desc)

	_, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		log.Printf("%v", err)
		return uuid.Nil
	}

	return task.ID
}

func (db *PostgresDatabase) Update(id, title, desc string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}

	var query string

	query = fmt.Sprintf(`update task set title = '%s', description = '%s' where id = '%s';`, title, desc, id)

	update, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error updating task with ID %s with error %v", id, err)
	}

	if update.RowsAffected() == 0 {
		return nil, fmt.Errorf("task with ID %v does not exist", UUID)
	}

	query = fmt.Sprintf("select * from tasks where id = '%s';", id)

	var updatedTask task.Task

	err = db.Conn.QueryRow(context.Background(), query).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Desc)
	if err != nil {
		return nil, fmt.Errorf("error updating task: %v", err)
	}

	return &updatedTask, nil
}

func (db *PostgresDatabase) Delete(id string) error {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %s", id)
	}

	query := fmt.Sprintf("delete from task where id = '%s';", id)

	del, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error deleting task with ID %v with error %v", id, err)
	}

	if del.RowsAffected() == 0 {
		return fmt.Errorf("task with ID %v does not exist", UUID)
	}

	log.Printf("Task with ID: %v has been deleted\n", UUID)
	return nil
}

func (db *PostgresDatabase) BulkImport(tasks []*task.Task) {
	for _, t := range tasks {
		query := fmt.Sprintf("insert into task (title, description) values ('%s', '%s');", t.Title, t.Desc)

		_, err := db.Conn.Exec(context.Background(), query)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}
