package repos

import (
	"Backend/internal/config"
	"Backend/internal/pkg/ctx_util"
	"Backend/internal/repos/scripts"
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error)
}

func (r *authRepo) StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error) {

	// 1. Define keys & args
	keys := []string{
		config.RedisKeyOTP(email, scene, codeID),
		config.RedisKeyThrottle(email, scene),
	}
	args := []interface{}{code, otpTTL, throttleTTL}

	// 2. Query
	res, err := r.scripts.StoreOTPAndThrottle.Run(ctx, r.rdb, keys, args...).Result()
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return false, err
		}
		return false, fmt.Errorf("StoreOTPAndThrottle lua error: %w", err)
	}

	// 3. Check reply
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 1 {
		return false, fmt.Errorf("unexpected lua reply: %v", res)
	}

	// 4. Check throttle
	status, _ := arr[0].(string)
	if status == "THROTTLED" {
		return true, nil
	}

	return false, nil
}

type authRepo struct {
	db      *sql.DB
	rdb     *redis.Client
	scripts *scripts.Registry
}

func NewAuthRepo(db *sql.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb, scripts: scripts.NewRegistry()}
}
