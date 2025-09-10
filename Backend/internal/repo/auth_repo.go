package repo

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type AuthRepo interface{}

type authRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAuthRepo(db *sql.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb}
}
