package repos

import (
	"backend/internal/config"
	"backend/internal/pkg/ctx_util"
	"backend/internal/repos/scripts"
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code string) error
	StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error)
}

// ThrottleMatchAndConsumeCode Check verify throttle -> Match code -> Consume code
func (r *authRepo) ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code string) error {
	return nil
}

// StoreOTPAndThrottle Check throttle -> Set throttle -> Store code
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
			return false, ctx.Err()
		}
		return false, ErrRunScript
	}

	// 3. Check reply
	arr, ok := res.([]interface{})
	if !ok || len(arr) != 1 {
		return false, ErrUnexpectedReply
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
