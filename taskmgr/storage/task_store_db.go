package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MrShanks/Taska/common/task"
)

type TaskStore struct {
	Conn *pgx.Conn
}

func (db *TaskStore) GetOne(id string) (*task.Task, error) {
	query := fmt.Sprintf("SELECT * FROM task WHERE id = '%s';", id)

	t := task.Task{}
	var author string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := db.Conn.QueryRow(ctx, query).Scan(&t.ID, &t.Title, &t.Desc, &author)
	if row == pgx.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}

	return &t, nil
}

func (db *TaskStore) GetTasks() []*task.Task {
	var fetchedTasks []*task.Task
	query := "SELECT * FROM task"

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rows, err := db.Conn.Query(ctx, query)
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

func (db *TaskStore) New(task *task.Task) uuid.UUID {
	// To be removed when proper user logic is implemented
	query := "SELECT id FROM author WHERE email = 'marco@rossi.com';"
	var authorID string

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := db.Conn.QueryRow(ctx, query).Scan(&authorID)
	if err == pgx.ErrNoRows {
		log.Printf("Couldn't find a match: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query = fmt.Sprintf("INSERT INTO task (title, description, author_id) VALUES ('%s', '%s', '%s');", task.Title, task.Desc, authorID)
	_, err = db.Conn.Exec(ctx, query)
	if err != nil {
		log.Printf("Could not insert new record into the database %v", err)
		return uuid.Nil
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query = fmt.Sprintf("SELECT id FROM task WHERE title = '%s'", task.Title)
  
	err = db.Conn.QueryRow(ctx, query).Scan(&task.ID)
	if err == pgx.ErrNoRows {
		log.Printf("Couldn't find a match: %v", err)
	}

	return task.ID
}

func (db *TaskStore) Update(id, title, desc string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}

	query := fmt.Sprintf(`UPDATE task SET title = '%s', description = '%s' WHERE id = '%s';`, title, desc, id)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	update, err := db.Conn.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error updating task with ID %s with error %v", id, err)
	}

	if update.RowsAffected() == 0 {
		return nil, fmt.Errorf("task with ID %v does not exist", UUID)
	}

	query = fmt.Sprintf("SELECT * FROM task WHERE id = '%s';", id)

	var updatedTask task.Task

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.Conn.QueryRow(ctx, query).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Desc, &updatedTask.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("error querying updated task: %v", err)
	}

	return &updatedTask, nil
}

func (db *TaskStore) Delete(id string) error {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %s", id)
	}

	query := fmt.Sprintf("DELETE FROM task WHERE id = '%s';", id)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	del, err := db.Conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error deleting task with ID %v with error %v", id, err)
	}

	if del.RowsAffected() == 0 {
		return fmt.Errorf("task with ID %v does not exist", UUID)
	}

	log.Printf("Task with ID: %v has been deleted\n", UUID)
	return nil
}

func (db *TaskStore) BulkImport(tasks []*task.Task) {
	for _, t := range tasks {
		query := fmt.Sprintf("INSERT INTO task (title, description) VALUES ('%s', '%s');", t.Title, t.Desc)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := db.Conn.Exec(ctx, query)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}
