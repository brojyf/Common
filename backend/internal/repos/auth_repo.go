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
	ThrottleLoginStoreDIDAndSession(ctx context.Context, email, password, deviceID string) error
	CheckOTTAndWriteUser(ctx context.Context, email, scene, jti, pwd string) error
	ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code, jti string, limit, window, jtiTTL int) error
	StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error)
}

// LoginStoreDIDAndSession TX: Store device id -> Store session
func (r *authRepo) ThrottleLoginStoreDIDAndSession(ctx context.Context, email, password, deviceID string) error {

	// 1. Run Lua

	// 2. Write SQL

	return nil
}

// CheckOTTAndWriteUser Check OTT -> Write user (Failed undo OTT)
func (r *authRepo) CheckOTTAndWriteUser(ctx context.Context, email, scene, jti, pwd string) error {

	// 1. Run Lua

	// 2. Write SQL (Failed: write back)

	return nil
}

// ThrottleMatchAndConsumeCode Check verify throttle -> Match code -> Consume code -> Set jti unused
func (r *authRepo) ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code, jti string, limit, window, jtiTTL int) error {

	// 1. Define keys & args
	keys := []string{
		config.RedisKeyVerifyThrottle(email, scene),
		config.RedisKeyOTP(email, scene, codeID),
		config.RedisKeyOTTJTIUsed(email, scene, jti),
	}
	args := []interface{}{limit, window, code, jtiTTL}

	// 2. Query
	res, err := r.scripts.ThrottleMatchAndConsumeCode.Run(ctx, r.rdb, keys, args...).Result()
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return ctx.Err()
		}
		return ErrRunScript
	}

	// 3, Check reply
	status, ok := res.(string)
	if !ok {
		return ErrUnexpectedReply
	}
	switch status {
	case "OK":
		return nil
	case "INVALID":
		return ErrOTPInvalid
	case "EXPIRED":
		return ErrOTPExpired
	case "THROTTLED":
		return ErrRateLimited
	default:
		return ErrUnexpectedReply
	}
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
	status, ok := res.(string)
	if !ok {
		return false, ErrUnexpectedReply
	}

	// 4. Check throttle
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
