package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kaanchinar/todo-app/models"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(addr, password string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		}),
		ttl: 5 * time.Minute,
	}
}

func (r *RedisStore) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}

func todoCacheKey(userID int) string {
	return fmt.Sprintf("todos:%d", userID)
}

func (r *RedisStore) InvalidateTodos(ctx context.Context, userID int) error {
	return r.client.Del(ctx, todoCacheKey(userID)).Err()
}

func (r *RedisStore) GetTodos(ctx context.Context, userID int) ([]models.Todo, bool) {
	data, err := r.client.Get(ctx, todoCacheKey(userID)).Bytes()
	if err != nil {
		return nil, false
	}

	var todos []models.Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, false
	}
	return todos, true
}

func (r *RedisStore) SetTodos(ctx context.Context, userID int, todos []models.Todo) error {
	data, err := json.Marshal(todos)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, todoCacheKey(userID), data, r.ttl).Err()
}
