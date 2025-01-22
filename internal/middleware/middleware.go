package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CustomContext struct {
	context.Context
	StartTime time.Time
}

type CustomHandler func(ctx *CustomContext, w http.ResponseWriter, r *http.Request)
type CustomMiddleware func(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error

// Wrap the CustomHandler to integrate with the http.Handler.
func AdaptHandler(handler CustomHandler, middleware ...CustomMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customContext := &CustomContext{
			Context:   context.Background(),
			StartTime: time.Now(),
		}
		//Run All middleware functions
		for _, mw := range middleware {
			err := mw(customContext, w, r)
			if err != nil {
				return
			}
		}

		fmt.Println("Run endpoint")
		handler(customContext, w, r)
		Log(customContext, w, r)
	}
}

func Log(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	elapsedTime := time.Since(ctx.StartTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
	return nil
}

func ParseForm(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	return nil
}

func ParseMultipartForm(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(10 << 20)
	return nil
}
