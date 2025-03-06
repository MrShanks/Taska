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
	// Tasks map[uuid.UUID]*task.Task
	Conn *pgx.Conn
}

func (db *PostgresDatabase) GetTasks() map[uuid.UUID]*task.Task {
	fetchedTasks := make(map[uuid.UUID]*task.Task)
	sqlString := "select * from tasks"
	rows, err := db.Conn.Query(context.Background(), sqlString)
	if err != nil {
		log.Printf("Error to query: %v", err)
	}

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
	sqlString := fmt.Sprintf("insert into tasks (id, title, description) values ('%s', '%s', '%s');", task.ID, task.Title, task.Desc)
	_, err := db.Conn.Exec(context.Background(), sqlString)
	if err != nil {
		log.Printf("%v", err)
	}
	return task.ID
}

func (db *PostgresDatabase) Delete(id string) error {
	UUID := uuid.MustParse(id)
	sqlString := fmt.Sprintf("select * from tasks where id = '%s';", id)
	db.Conn.QueryRow(context.Background(), sqlString)
	// if row == errors {
	// 	return fmt.Errorf("task with ID: %v not found", UUID)
	// }

	sqlString = fmt.Sprintf("delete from tasks where id = '%s';", id)
	_, err := db.Conn.Exec(context.Background(), sqlString)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("Task with ID: %v has been deleted", UUID)
	return nil
}
