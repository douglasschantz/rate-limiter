package entity

import (
	"errors"
	"time"
)

type RequestIPInputDTO struct {
	IP string
}

var (
	ErrKeyNotFound = errors.New("key not found")
)

type RequestRepositoryInterface interface {
	SaveRequestIP(respIp *Ip) error
	GetCountLimiter(respIp *Ip) (int, error)
	SaveToken(token string, expiration time.Duration) error
	GetToken(token string) (bool, error)
}
