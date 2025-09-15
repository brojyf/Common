package repo

import (
	"Backend/internal/config"
	"Backend/internal/x"
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	MatchAndConsumeOTP(ctx context.Context, email, scene, code string) (bool, error)
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

func (r *authRepo) MatchAndConsumeOTP(ctx context.Context, email, scene, code string) (bool, error) {
	pattern := fmt.Sprintf("otp:%s:%s:*", email, scene)

	iter := r.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if x.IsCtxDone(ctx, nil) {
			return false, ctx.Err()
		}
		k := iter.Val()

		v, err := r.rdb.Get(ctx, k).Result()
		if errors.Is(err, redis.Nil) {
			continue
		}
		if err != nil {
			return false, err
		}

		if subtle.ConstantTimeCompare([]byte(v), []byte(code)) == 1 {
			if err := r.rdb.Del(ctx, k).Err(); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	if err := iter.Err(); err != nil {
		return false, err
	}
	return false, nil
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
