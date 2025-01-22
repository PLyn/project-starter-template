package handlers

import (
	"cloud/internal/middleware"
	"cloud/internal/template/home"
	"net/http"
)

// Example function using CustomHandler type.
func MyHandler(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	view(w, r)
}

func view(w http.ResponseWriter, r *http.Request) {
	home.Template("Templ Quickstart").Render(r.Context(), w)
}
