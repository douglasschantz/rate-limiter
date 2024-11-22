package database

import (
	"context"
	"time"

	"github.com/douglasschantz/rate-limiter/internal/entity"
	"github.com/redis/go-redis/v9"
)

type RequestRepository struct {
	Redis *redis.Client
}

func NewRequestRepository(redis *redis.Client) *RequestRepository {
	return &RequestRepository{Redis: redis}

}

func (r *RequestRepository) SaveRequestIP(request *entity.Ip) error {
	var limite int = 2
	ctx := context.Background()
	_, err := r.Redis.HIncrBy(ctx, "requests", request.IP, 1).Result()
	if err != nil {
		return err
	}

	// Definir um tempo de expiração para o hash (1 segundo)
	_, err = r.Redis.Expire(ctx, "requests", time.Duration(limite)*time.Second).Result()
	if err != nil {

		return err
	}
	return nil
}

func (r *RequestRepository) GetCountLimiter(request *entity.Ip) (int, error) {
	ctx := context.Background()
	count, err := r.Redis.HGet(ctx, "requests", request.IP).Int()
	if err != nil && err != redis.Nil {
		return 0, err
	}

	return count, nil
}

// Começa apenas metodos do Token

func (r *RequestRepository) SaveToken(token string, expiration time.Duration) error {
	ctx := context.Background()
	// Salva o token no Redis
	err := r.Redis.Set(ctx, token, true, expiration).Err()
	if err != nil {
		return err
	}
	// Define o tempo de expiração
	err = r.Redis.Expire(ctx, token, expiration).Err()
	if err != nil {
		return err
	}
	return nil

}

func (r *RequestRepository) GetToken(token string) (bool, error) {
	ctx := context.Background()
	result, err := r.Redis.Get(ctx, token).Bool()
	if err != nil && err != redis.Nil {
		return false, err
	}
	return result, nil
}
