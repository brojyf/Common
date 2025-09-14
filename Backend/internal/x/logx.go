package x

import (
	"log"
)

// Error 打错误日志
func Error(rid, op string, err error) {
	log.Printf("[ERROR] rid=%s op=%s err=%v", rid, op, err)
}

// Info 打普通信息日志
func Info(rid, op, msg string) {
	log.Printf("[INFO] rid=%s op=%s msg=%s", rid, op, msg)
}
