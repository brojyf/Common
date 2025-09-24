package repos

import (
	"backend/internal/config"
	"backend/internal/pkg/ctx_util"
	"backend/internal/repos/scripts"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	StoreDIDAndSession(ctx context.Context, userID uint64, deviceID []byte, pushToken *string, rtkHash []byte, tokenVersion uint, expiresAt time.Time) error
	UndoOTTMark(email, scene, jti string, ttlSec int)
	CreateUser(ctx context.Context, email, pwd string) (uint64, uint, error)
	ConsumeOTTJTI(ctx context.Context, email, scene, jti string, newTTL int) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	ThrottleMatchAndConsumeCode(ctx context.Context, email, scene, codeID, code, jti string, limit, window, jtiTTL int) error
	StoreOTPAndThrottle(ctx context.Context, email, scene, codeID, code string, otpTTL, throttleTTL int) (bool, error)
}

func (r *authRepo) StoreDIDAndSession(
	ctx context.Context,
	userID uint64,
	deviceID []byte,
	pushToken *string,
	rtkHash []byte,
	tokenVersion uint,
	expiresAt time.Time) error {

	// 1. Check input
	if len(deviceID) != 16 {
		return fmt.Errorf("%w: invalid device_id length=%d (want 16)", ErrUnexpectedSQL, len(deviceID))
	}
	if len(rtkHash) != 32 {
		return fmt.Errorf("%w: invalid rtk_hash length=%d (want 32)", ErrUnexpectedSQL, len(rtkHash))
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%w: begin tx: %v", ErrUnexpectedSQL, err)
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Upsert user_devices
	const upsertDev = `
		INSERT INTO user_devices (device_id, user_id, push_token, last_seen_at)
		VALUES (?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
  			push_token   = VALUES(push_token),
  			last_seen_at = VALUES(last_seen_at),
  			revoked_at   = NULL
	`
	if _, err := tx.ExecContext(ctx, upsertDev, deviceID, userID, pushToken); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			switch me.Number {
			case 1452:
				return ErrNotFound
			}
		}
		return fmt.Errorf("%w: upsert user_devices: %v", ErrUnexpectedSQL, err)
	}

	// 2) Upsert session
	const upsertSess = `
		INSERT INTO sessions (user_id, device_id, rtk_hash, token_version, expires_at)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
  			rtk_hash      = VALUES(rtk_hash),
  			token_version = VALUES(token_version),
  			expires_at    = VALUES(expires_at),
  			revoked_at    = NULL
	`
	if _, err := tx.ExecContext(ctx, upsertSess, userID, deviceID, rtkHash, tokenVersion, expiresAt.UTC()); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) {
			switch me.Number {
			case 1452:
				return ErrNotFound
			}
		}
		return fmt.Errorf("%w: upsert sessions: %v", ErrUnexpectedSQL, err)
	}

	if err := tx.Commit(); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%w: commit: %v", ErrUnexpectedSQL, err)
	}
	return nil
}

// UndoOTTMark Undo one time token mark after sql failed
func (r *authRepo) UndoOTTMark(email, scene, jti string, ttlSec int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ttl := time.Duration(ttlSec) * time.Second
	_ = r.rdb.Set(ctx, config.RedisKeyOTTJTIUsed(email, scene, jti), 0, ttl).Err()
}

// CreateUser Write new user info to sql
func (r *authRepo) CreateUser(ctx context.Context, email, pwd string) (uint64, uint, error) {
	const query = `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, email, pwd)
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return 0, 0, ctx.Err()
		}
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 || strings.Contains(err.Error(), "Duplicate entry") {
			return 0, 0, ErrEmailAlreadyExists
		}
		return 0, 0, fmt.Errorf("%w:%s", ErrUnexpectedSQL, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, 0, fmt.Errorf("%w:%s", ErrUnexpectedSQL, err)
	}

	return uint64(id), 1, nil
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
		return fmt.Errorf("%w:%s", ErrRunScript, err)
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
		return false, fmt.Errorf("%w:%s", ErrUnexpectedSQL, err)
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
		return fmt.Errorf("%w:%s", ErrRunScript, err)
	}

	// 3. Check reply
	status, ok := res.(string)
	if !ok {
		return fmt.Errorf("%w:%s", ErrUnexpectedReply, err)
	}
	switch status {
	case "OK":
		return nil
	case "INVALID":
		return fmt.Errorf("%w:%s", ErrOTPInvalid, err)
	case "EXPIRED":
		return fmt.Errorf("%w:%s", ErrOTPExpired, err)
	case "THROTTLED":
		return fmt.Errorf("%w:%s", ErrRateLimited, err)
	default:
		return fmt.Errorf("%w:%s", ErrUnexpectedReply, err)
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
		return false, fmt.Errorf("%w:%s", ErrRunScript, err)
	}

	// 3. Check reply
	status, ok := res.(string)
	if !ok {
		return false, fmt.Errorf("%w:%s", ErrUnexpectedReply, err)
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
