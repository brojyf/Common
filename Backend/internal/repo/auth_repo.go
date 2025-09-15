package repo

import (
	"Backend/internal/config"
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	StoreCode(ctx context.Context, code, email, scene string) error
	SetThrottle(ctx context.Context, email, scene string) error
}

type authRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAuthRepo(db *sql.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb}
}

func (a *authRepo) StoreCode(ctx context.Context, code, email, scene string) error {
	k := config.RedisKeyOTP(email, scene)
	return a.rdb.Set(ctx, k, code, config.OTP_TTL).Err()
}

func (a *authRepo) SetThrottle(ctx context.Context, email, scene string) error {
	key := config.RedisKeyThrottle(email, scene)
	return a.rdb.Set(ctx, key, "1", config.OTP_THROTTLE_TTL).Err()
}
