package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sky0621/techcv-app/backend/internal/shared/httpserver"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: httpserver.NewRouter(),
	}

	log.Printf("server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
