package logx

import (
	"Backend/internal/pkg/request_id"
	"context"
	"log"
)

// LogError Print error log
func LogError(c context.Context, op string, err error) {
	rid, _ := request_id.From(c)
	log.Printf("\033[31m[ERROR]\033[0m request_id=%s op=%s err=%v", rid, op, err)
}

// LogInfo Print information log
func LogInfo(c context.Context, op, msg string) {
	rid, _ := request_id.From(c)
	log.Printf("[INFO] request_id=%s op=%s msg=%s", rid, op, msg)
}
