package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kaanchinar/todo-app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	s := &PostgresStore{db: pool}
	if err := s.migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return s, nil
}

func (s *PostgresStore) migrate(ctx context.Context) error {
	_, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
	`)
	return err
}

func (s *PostgresStore) Close() {
	s.db.Close()
}

// User methods
func (s *PostgresStore) CreateUser(ctx context.Context, user *models.User) error {
	err := s.db.QueryRow(ctx,
		"INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, time.Now().UTC(),
	).Scan(&user.ID)
	return err
}

func (s *PostgresStore) GetUserByUsername(ctx context.Context, username string) (*models.User, bool) {
	var user models.User
	err := s.db.QueryRow(ctx,
		"SELECT id, username, password, created_at FROM users WHERE username = $1", username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (s *PostgresStore) GetUserByID(ctx context.Context, id int) (*models.User, bool) {
	var user models.User
	err := s.db.QueryRow(ctx,
		"SELECT id, username, password, created_at FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return &user, true
}

// Todo methods
func (s *PostgresStore) CreateTodo(ctx context.Context, todo *models.Todo) error {
	now := time.Now().UTC()
	err := s.db.QueryRow(ctx,
		"INSERT INTO todos (title, completed, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		todo.Title, todo.Completed, todo.UserID, now, now,
	).Scan(&todo.ID)
	if err != nil {
		return err
	}
	todo.CreatedAt = now
	todo.UpdatedAt = now
	return nil
}

func (s *PostgresStore) GetTodosByUserID(ctx context.Context, userID int) ([]models.Todo, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, title, completed, user_id, created_at, updated_at FROM todos WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.UserID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *PostgresStore) GetTodoByID(ctx context.Context, id int) (*models.Todo, bool) {
	var todo models.Todo
	err := s.db.QueryRow(ctx,
		"SELECT id, title, completed, user_id, created_at, updated_at FROM todos WHERE id = $1", id,
	).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return &todo, true
}

func (s *PostgresStore) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	_, err := s.db.Exec(ctx,
		"UPDATE todos SET title = $1, completed = $2, updated_at = $3 WHERE id = $4",
		todo.Title, todo.Completed, time.Now().UTC(), todo.ID,
	)
	return err
}

func (s *PostgresStore) DeleteTodo(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}
