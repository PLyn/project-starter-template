package main

import (
	"cloud/internal/SSE"
	"cloud/internal/generate"
	"cloud/internal/middleware"
	"cloud/internal/template/backup"
	"cloud/internal/template/home"
	"cloud/internal/view"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := generate.GenerateMain()
	if err != nil {
		panic(err)
	}

	_ = godotenv.Load()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	mux.HandleFunc("GET /static/", view.ServeStaticFiles)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		middleware.Chain(w, r, home.Template("Templ Quickstart"))
	})

	mux.HandleFunc("GET /backup", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/backup" {
			http.NotFound(w, r)
			return
		}
		middleware.Chain(w, r, backup.Template("Backup"))
	})

	mux.HandleFunc("GET /feed", SSE.Feed)

	fmt.Printf("server is running on port %s\n", os.Getenv("PORT"))
	err = http.ListenAndServe("127.0.0.1:"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}
}
