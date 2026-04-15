package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlist-service/internal/bootstrap"
)

func main() {
	ctx := context.Background()
	app, err := bootstrap.NewApp(ctx)
	if err != nil {
		log.Fatalf("bootstrap app: %v", err)
	}
	defer app.Close()
	srv := app.Server()

	go func() {
		log.Printf("HTTP server started at %s", app.HTTPAddr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
