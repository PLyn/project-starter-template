package middleware

import (
	"net/http"
	"time"
)

func authorize(ctx *CustomContext, _ http.ResponseWriter, _ *http.Request) error {
	currentTime := time.Now().UTC()
	ctx.StartTime = currentTime // set your new context value
	return nil
}
