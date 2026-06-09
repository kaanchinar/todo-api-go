package store

import (
	"context"

	"github.com/kaanchinar/todo-app/models"
)

type Store interface {
	// Users
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, bool)
	GetUserByID(ctx context.Context, id int) (*models.User, bool)

	// Todos
	CreateTodo(ctx context.Context, todo *models.Todo) error
	GetTodosByUserID(ctx context.Context, userID int) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id int) (*models.Todo, bool)
	UpdateTodo(ctx context.Context, todo *models.Todo) error
	DeleteTodo(ctx context.Context, id int) error

	// Cache
	InvalidateTodos(ctx context.Context, userID int) error
	GetTodos(ctx context.Context, userID int) ([]models.Todo, bool)
	SetTodos(ctx context.Context, userID int, todos []models.Todo) error

	// Lifecycle
	Close()
}

type CompositeStore struct {
	postgres *PostgresStore
	redis    *RedisStore
}

func NewCompositeStore(ctx context.Context, databaseURL, redisAddr, redisPass string) (*CompositeStore, error) {
	pg, err := NewPostgresStore(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	r := NewRedisStore(redisAddr, redisPass)
	if err := r.Ping(ctx); err != nil {
		pg.Close()
		return nil, err
	}

	return &CompositeStore{postgres: pg, redis: r}, nil
}

func (c *CompositeStore) Close() {
	c.postgres.Close()
	c.redis.Close()
}

// User methods — always go to Postgres
func (c *CompositeStore) CreateUser(ctx context.Context, user *models.User) error {
	return c.postgres.CreateUser(ctx, user)
}

func (c *CompositeStore) GetUserByUsername(ctx context.Context, username string) (*models.User, bool) {
	return c.postgres.GetUserByUsername(ctx, username)
}

func (c *CompositeStore) GetUserByID(ctx context.Context, id int) (*models.User, bool) {
	return c.postgres.GetUserByID(ctx, id)
}

// Todo methods — Postgres for writes, Redis cache for reads
func (c *CompositeStore) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if err := c.postgres.CreateTodo(ctx, todo); err != nil {
		return err
	}
	c.redis.InvalidateTodos(ctx, todo.UserID)
	return nil
}

func (c *CompositeStore) GetTodosByUserID(ctx context.Context, userID int) ([]models.Todo, error) {
	if cached, ok := c.redis.GetTodos(ctx, userID); ok {
		return cached, nil
	}

	todos, err := c.postgres.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	c.redis.SetTodos(ctx, userID, todos)
	return todos, nil
}

func (c *CompositeStore) GetTodoByID(ctx context.Context, id int) (*models.Todo, bool) {
	return c.postgres.GetTodoByID(ctx, id)
}

func (c *CompositeStore) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	if err := c.postgres.UpdateTodo(ctx, todo); err != nil {
		return err
	}
	c.redis.InvalidateTodos(ctx, todo.UserID)
	return nil
}

func (c *CompositeStore) DeleteTodo(ctx context.Context, id int) error {
	todo, ok := c.postgres.GetTodoByID(ctx, id)
	if !ok {
		return nil
	}
	if err := c.postgres.DeleteTodo(ctx, id); err != nil {
		return err
	}
	c.redis.InvalidateTodos(ctx, todo.UserID)
	return nil
}

// Cache passthrough
func (c *CompositeStore) InvalidateTodos(ctx context.Context, userID int) error {
	return c.redis.InvalidateTodos(ctx, userID)
}

func (c *CompositeStore) GetTodos(ctx context.Context, userID int) ([]models.Todo, bool) {
	return c.redis.GetTodos(ctx, userID)
}

func (c *CompositeStore) SetTodos(ctx context.Context, userID int, todos []models.Todo) error {
	return c.redis.SetTodos(ctx, userID, todos)
}
