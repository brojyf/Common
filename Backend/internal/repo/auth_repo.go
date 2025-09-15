package repo

import (
	"Backend/internal/config"
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	MatchAndConsumeOTP(ctx context.Context, email, scene, code, jti string) (bool, error)
	CheckThrottle(ctx context.Context, email, scene string) (bool, error)
	StoreCode(ctx context.Context, code, email, scene, jti string) error
	SetThrottle(ctx context.Context, email, scene string) error
}

type authRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAuthRepo(db *sql.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb}
}

func (r *authRepo) MatchAndConsumeOTP(ctx context.Context, email, scene, code, jti string) (bool, error) {
	key := fmt.Sprintf("otp:%s:%s:%s", email, scene, jti)
	v, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(v), []byte(code)) != 1 {
		return false, nil
	}
	got, err := r.rdb.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if subtle.ConstantTimeCompare([]byte(got), []byte(code)) != 1 {
		return false, nil
	}
	return true, nil
}

func (r *authRepo) CheckThrottle(ctx context.Context, email, scene string) (bool, error) {
	key := config.RedisKeyThrottle(email, scene)
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *authRepo) StoreCode(ctx context.Context, code, email, scene, jti string) error {
	k := config.RedisKeyOTP(email, scene, jti)
	return r.rdb.Set(ctx, k, code, config.C.RedisTTL.OTP).Err()
}

func (r *authRepo) SetThrottle(ctx context.Context, email, scene string) error {
	key := config.RedisKeyThrottle(email, scene)
	return r.rdb.Set(ctx, key, "1", config.C.RedisTTL.OTPThrottle).Err()
}
