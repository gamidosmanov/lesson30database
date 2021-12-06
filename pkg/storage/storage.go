package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database abstraction
type Storage struct {
	db *pgxpool.Pool
}

// New returns Storage object with pg connection pool
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{db: db}
	return &s, nil
}

// Task represents single task
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks selects task list from database
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM devbase.tasks.tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
		`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// TasksByLabel returns all tasks by given label
func (s *Storage) TasksByLabel(labelID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM
			tasks.tasks,
			tasks.tasks_labels
		WHERE
			tasks.tasks.id = tasks.tasks_labels.task_id
			AND tasks.tasks_labels.label_id = $1
		;
		`,
		labelID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// NewTask
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO devbase.tasks.tasks (opened, closed, author_id, assigned_id, title, content)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
		`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// DeleteTask
func (s *Storage) DeleteTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		DELETE FROM devbase.tasks.tasks
		WHERE id = $1 RETURNING id
		`,
		t.ID,
	).Scan(&id)
	return id, err
}

// UpdateTask
func (s *Storage) UpdateTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		UPDATE devbase.tasks.tasks
		SET opened = $2
			, closed = $3
			, author_id = $4
			, assigned_id = $5
			, title = $6
			, content = $7
		WHERE id = $1
		RETURNING id
		`,
		t.ID,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}
