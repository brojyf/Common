package x

import (
	"Backend/internal/pkg/request_id"
	"log"

	"golang.org/x/net/context"
)

// Error Print error log
func Error(c context.Context, op string, err error) {
	rid, _ := request_id.From(c)
	log.Printf("[ERROR] request_id=%s op=%s err=%v", rid, op, err)
}

// Info Print information log
func Info(c context.Context, op, msg string) {
	rid, _ := request_id.From(c)
	log.Printf("[INFO] request_id=%s op=%s msg=%s", rid, op, msg)
}
