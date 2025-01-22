package main

import (
	"cloud/cmd/handlers"
	"cloud/cmd/sseHandlers"
	"cloud/internal/middleware"
	"cloud/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", utils.ServeFavicon)
	mux.HandleFunc("GET /static/", utils.ServeStaticFiles)

	mux.Handle("/backup", middleware.AdaptHandler(handlers.Backup))
	mux.HandleFunc("GET /backupsse", sseHandlers.Backup)

	fmt.Printf("server is running on port %s\n", os.Getenv("PORT"))
	err := http.ListenAndServe("127.0.0.1:"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}
}
