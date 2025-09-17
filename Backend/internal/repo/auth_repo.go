package repo

import (
	"Backend/internal/config"
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	UpdateUsername(ctx context.Context, uid uint64, username string) error
	StoreSession(ctx context.Context, uid uint64, deviceID, rtkHashed []byte, tkv uint, exp time.Time) error
	GetTKVersion(ctx context.Context, uid uint64) (uint, error)
	StoreDeviceID(ctx context.Context, uid uint64, deviceID []byte) error
	StoreNewUser(ctx context.Context, email, hashPwd string) (uint64, error)
	// Redis Methods
	MatchAndConsumeOTP(ctx context.Context, email, scene, code, jti string) (bool, error)
	StoreCode(ctx context.Context, code, email, scene, jti string) error
	CheckThrottle(ctx context.Context, email, scene string) (bool, error)
	SetThrottle(ctx context.Context, email, scene string) error
}

type authRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAuthRepo(db *sql.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb}
}

func (r *authRepo) UpdateUsername(ctx context.Context, uid uint64, username string) error {
	const q = `
	UPDATE users
   	  SET username   = ?,
       updated_at = CURRENT_TIMESTAMP
 	  WHERE id = ?`

	res, err := r.db.ExecContext(ctx, q, username, uid)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return nil
		}
		return err
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *authRepo) StoreSession(ctx context.Context, uid uint64, deviceID, rtkHashed []byte, tkv uint, exp time.Time) error {
	if len(deviceID) != 16 {
		return fmt.Errorf("invalid device id")
	}

	const q = `
INSERT INTO sessions (user_id, device_id, rtk_hash, token_version, expires_at, revoked_at)
VALUES (?, ?, ?, ?, ?, NULL)
ON DUPLICATE KEY UPDATE
  rtk_hash      = VALUES(rtk_hash),
  token_version = VALUES(token_version),
  expires_at    = VALUES(expires_at),
  revoked_at    = NULL,
  updated_at    = CURRENT_TIMESTAMP
`
	_, err := r.db.ExecContext(ctx, q, uid, deviceID, rtkHashed, tkv, exp)
	return err
}

func (r *authRepo) GetTKVersion(ctx context.Context, uid uint64) (uint, error) {
	const q = `SELECT token_version FROM users WHERE id = ?`
	var v int
	err := r.db.QueryRowContext(ctx, q, uid).Scan(&v)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}
	return uint(v), nil
}

func (r *authRepo) StoreDeviceID(ctx context.Context, uid uint64, deviceID []byte) error {
	const q = `
	INSERT INTO user_devices (device_id, user_id, revoked_at)
	VALUES (?, ?, NULL)
	ON DUPLICATE KEY UPDATE
  	  	user_id = VALUES(user_id),
  	  	revoked_at = NULL
	`
	_, err := r.db.ExecContext(ctx, q, deviceID, uid)
	return err
}

// StoreNewUser (0, nil) represents conflict
func (r *authRepo) StoreNewUser(ctx context.Context, email, hashPwd string) (uint64, error) {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`

	res, err := r.db.ExecContext(ctx, query, email, hashPwd)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return 0, nil
		}
		return 0, err
	} // Check conflict

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	} // Check get uid
	return uint64(id), nil
}

func (r *authRepo) MatchAndConsumeOTP(ctx context.Context, email, scene, code, codeID string) (bool, error) {
	key := fmt.Sprintf("otp:%s:%s:%s", email, scene, codeID)
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

func (r *authRepo) StoreCode(ctx context.Context, code, email, scene, id string) error {
	k := config.RedisKeyOTP(email, scene, id)
	return r.rdb.Set(ctx, k, code, config.C.RedisTTL.OTP).Err()
}

func (r *authRepo) SetThrottle(ctx context.Context, email, scene string) error {
	key := config.RedisKeyThrottle(email, scene)
	return r.rdb.Set(ctx, key, "1", config.C.RedisTTL.OTPThrottle).Err()
}
