package repo

import (
	"Backend/internal/config"
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	CheckThrottle(ctx context.Context, email, scene string) (bool, error)
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

func (r *authRepo) CheckThrottle(ctx context.Context, email, scene string) (bool, error) {
	key := config.RedisKeyThrottle(email, scene)
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *authRepo) StoreCode(ctx context.Context, code, email, scene string) error {
	k := config.RedisKeyOTP(email, scene)
	return r.rdb.Set(ctx, k, code, config.C.RedisTTL.OTP).Err()
}

func (r *authRepo) SetThrottle(ctx context.Context, email, scene string) error {
	key := config.RedisKeyThrottle(email, scene)
	return r.rdb.Set(ctx, key, "1", config.C.RedisTTL.OTPThrottle).Err()
}
