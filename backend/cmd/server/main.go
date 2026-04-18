package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/sky0621/techcv-app/backend/internal/shared/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	application, err := app.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("close app: %v", err)
		}
	}()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: application.Handler,
	}

	log.Printf("server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
