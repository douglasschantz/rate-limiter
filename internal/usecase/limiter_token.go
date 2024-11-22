package usecase

import (
	"context"

	"github.com/douglasschantz/rate-limiter/internal/entity"
)

type RateLimiterByToken struct {
	redisRepository entity.RequestRepositoryInterface
}

func NewRateLimiterByToken(redisRepository entity.RequestRepositoryInterface) *RateLimiterByToken {
	return &RateLimiterByToken{redisRepository: redisRepository}
}
func (rl *RateLimiterByToken) CheckRateLimitByToken(ctx context.Context, token string) (bool, error) {
	key := token
	// Obtém o contador atual do Redis
	result, err := rl.redisRepository.GetToken(key)
	if err != nil {
		if err != entity.ErrKeyNotFound {
			return false, err
		}
		return false, nil
	}

	return result, nil
}
