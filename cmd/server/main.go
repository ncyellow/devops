package main

import (
	"context"
	"github.com/ncyellow/devops/internal/server/handlers"
	"github.com/ncyellow/devops/internal/server/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	repo := storage.NewRepository()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	http.HandleFunc("/update/gauge/", handlers.Handler(repo))
	http.HandleFunc("/update/counter/", handlers.Handler(repo))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
