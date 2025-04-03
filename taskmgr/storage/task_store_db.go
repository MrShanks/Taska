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

func (db *TaskStore) Search(keyword, authorID string) ([]*task.Task, error) {
	var fetchedTasks []*task.Task
	query := fmt.Sprintf("SELECT * FROM task WHERE author_id = '%s' AND (title LIKE '%%%s%%' OR description LIKE '%%%s%%');", authorID, keyword, keyword)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rows, err := db.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		r := &task.Task{}

		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.AuthorID, &r.Complete)
		if err != nil {
			return nil, err
		}

		fetchedTasks = append(fetchedTasks, r)
	}

	return fetchedTasks, nil
}

func (db *TaskStore) GetOne(id, authorID string) (*task.Task, error) {
	query := fmt.Sprintf("SELECT * FROM task WHERE id = '%s';", id)

	t := task.Task{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row := db.Conn.QueryRow(ctx, query).Scan(&t.ID, &t.Title, &t.Desc, &t.AuthorID, &t.Complete)
	if row == pgx.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}

	return &t, nil
}

func (db *TaskStore) GetTasks(authorID string) []*task.Task {
	var fetchedTasks []*task.Task
	query := fmt.Sprintf("SELECT * FROM task WHERE author_id = '%s'", authorID)

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

		err := rows.Scan(&r.ID, &r.Title, &r.Desc, &r.AuthorID, &r.Complete)
		if err != nil {
			log.Printf("%v", err)
		}

		fetchedTasks = append(fetchedTasks, r)
	}
	return fetchedTasks
}

func (db *TaskStore) New(task *task.Task) uuid.UUID {
	query := fmt.Sprintf("INSERT INTO task (author_id, title, description) VALUES ('%s', '%s', '%s');", task.AuthorID, task.Title, task.Desc)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := db.Conn.Exec(ctx, query)
	if err != nil {
		log.Printf("Could not insert new record into the database %v", err)
		return uuid.Nil
	}

	query = fmt.Sprintf("SELECT id FROM task WHERE title='%s'", task.Title)
	row := db.Conn.QueryRow(ctx, query).Scan(&task.ID)
	if row == pgx.ErrNoRows {
		log.Printf("No task was found with ID: %v", err)
		return uuid.Nil
	}

	return task.ID
}

func (db *TaskStore) Update(id, title, desc, authorID string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}

	query := fmt.Sprintf(`UPDATE task SET title = '%s', description = '%s' WHERE id = '%s' and author_id = '%s';`, title, desc, id, authorID)

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

	err = db.Conn.QueryRow(ctx, query).Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Desc, &updatedTask.AuthorID, &updatedTask.Complete)
	if err != nil {
		return nil, fmt.Errorf("error querying updated task: %v", err)
	}

	return &updatedTask, nil
}

func (db *TaskStore) Delete(id, authorID string) error {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %s", id)
	}

	query := fmt.Sprintf("DELETE FROM task WHERE id = '%s' and author_id = '%s';", id, authorID)

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

func (db *TaskStore) BulkImport(tasks []*task.Task, authorID string) {
	for _, t := range tasks {
		query := fmt.Sprintf("INSERT INTO task (title, description, author_id) VALUES ('%s', '%s', '%s');", t.Title, t.Desc, authorID)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := db.Conn.Exec(ctx, query)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}
