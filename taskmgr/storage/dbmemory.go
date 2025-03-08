package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MrShanks/Taska/common/task"
)

type PostgresDatabase struct {
	Conn *pgx.Conn
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

func (db *PostgresDatabase) Delete(id string) error {
	UUID := uuid.MustParse(id)
	query := fmt.Sprintf("delete from tasks where id = '%s';", id)

	del, err := db.Conn.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Error deleting task with ID %v: %v\n", id, err)
		return err
	} else if del.String() == "DELETE 0" {
		log.Printf("Task with ID %v does not exist\n", UUID)
		return fmt.Errorf("task with ID %v does not exist", UUID)
	}
	log.Printf("Task with ID: %v has been deleted\n", UUID)
	return nil
}
