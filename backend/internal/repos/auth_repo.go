package repos

import (
	"backend/internal/config"
	"backend/internal/pkg/ctx_util"
	"backend/internal/repos/scripts"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	StoreDIDAndSession(ctx context.Context, uid uint64, did string) error
	UndoOTTMark(ctx context.Context, email, scene, jti string, ttlSec int)
	CreateUser(ctx context.Context, email, pwd string) (uint64, error)
	ConsumeOTTJTI(ctx context.Context, email, scene, jti string, newTTL int) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code, jti string, limit, window, jtiTTL int) error
	StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error)
}

func (r *authRepo) StoreDIDAndSession(ctx context.Context, uid uint64, did string) error {

	return nil
}

// UndoOTTMark Undo one time token mark after sql failed
func (r *authRepo) UndoOTTMark(ctx context.Context, email, scene, jti string, ttlSec int) {
	ttl := time.Duration(ttlSec) * time.Second
	_ = r.rdb.Set(ctx, config.RedisKeyOTTJTIUsed(email, scene, jti), 0, ttl).Err()
}

// CreateUser Write new user info to sql
func (r *authRepo) CreateUser(ctx context.Context, email, pwd string) (uint64, error) {
	const query = `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, email, pwd)
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return 0, ctx.Err()
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 || strings.Contains(err.Error(), "Duplicate entry") {
			return 0, ErrEmailAlreadyExists
		}
		return 0, ErrUnexpectedSQL
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, ErrUnexpectedSQL
	}

	return uint64(id), nil
}

// ConsumeOTTJTI Consume one time token jti
func (r *authRepo) ConsumeOTTJTI(ctx context.Context, email, scene, jti string, newTTL int) error {
	keyStr := config.RedisKeyOTTJTIUsed(email, scene, jti)
	key := []string{keyStr}
	arg := []interface{}{newTTL}
	cmd := r.scripts.FindAdnMarkOTTJTI.Run(ctx, r.rdb, key, arg...)
	n, err := cmd.Int()
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return ctx.Err()
		}
		return ErrRunScript
	}
	if n != 0 {
		return ErrOTPInvalid
	}
	return nil
}

// CheckEmailExists Check if email exists
func (r *authRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var dummy int
	err := r.db.QueryRowContext(ctx,
		"SELECT 1 FROM users WHERE email = ? LIMIT 1", email,
	).Scan(&dummy)
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return false, ctx.Err()
		}
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, ErrUnexpectedSQL
	}
	return true, nil
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

	// 3. Check reply
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
