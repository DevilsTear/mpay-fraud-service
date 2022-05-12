package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Repository represent the repositories
type Repository interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
}

// repository represent the repository model
type repository struct {
	Client  redis.Cmdable
	Context context.Context
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(Client redis.Cmdable) Repository {
	return &repository{Client: Client}
}

// Set attaches the redis repository and set the data
func (r *repository) Set(key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(r.Context, key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(key string) (string, error) {
	get := r.Client.Get(r.Context, key)
	return get.Result()
}
