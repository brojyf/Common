package models

type User struct {
	UserID  uint64
	Email   string
	PwdHash string
	TokenV  uint
}
