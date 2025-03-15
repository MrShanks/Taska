package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MrShanks/Taska/common/task"
)

type TaskStoreDB struct {
	Conn *pgx.Conn
}

func (db *TaskStoreDB) GetOne(id string) (*task.Task, error) {
	query := fmt.Sprintf("select * from task where id = '%s';", id)

	t := task.Task{}
	var author string

	row := db.Conn.QueryRow(context.Background(), query).Scan(&t.ID, &t.Title, &t.Desc, &author)
	if row == pgx.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}

	return &t, nil
}

func (db *TaskStoreDB) GetTasks() []*task.Task {
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
		var author string

		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &author)
		if err != nil {
			log.Printf("%v", err)
		}

		fetchedTasks = append(fetchedTasks, r)
	}
	return fetchedTasks
}

func (db *TaskStoreDB) New(task *task.Task) uuid.UUID {
	// To be removed when proper user logic is implemented
	queryID := "SELECT id FROM author WHERE email = 'marco@rossi.com';"
	var ID string

	err := db.Conn.QueryRow(context.Background(), queryID).Scan(&ID)
	if err == pgx.ErrNoRows {
		log.Printf("Couldn't find a match: %v", err)
	}

	query := fmt.Sprintf("insert into task (author_id, title, description) values ('%s', '%s', '%s');", ID, task.Title, task.Desc)
	_, err = db.Conn.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Could not insert new record into the database %v", err)
		return uuid.Nil
	}

	return task.ID
}

func (db *TaskStoreDB) Update(id, title, desc string) (*task.Task, error) {
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

func (db *TaskStoreDB) Delete(id string) error {
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

func (db *TaskStoreDB) BulkImport(tasks []*task.Task) {
	for _, t := range tasks {
		query := fmt.Sprintf("insert into tasks (title, description) values ('%s', '%s');", t.Title, t.Desc)

		_, err := db.Conn.Exec(context.Background(), query)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}
