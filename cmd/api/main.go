package main

import (
	"context"
	"fmt"
	"log"
	"mekoko/internal/config"
	"mekoko/internal/server"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func run(ctx context.Context) error {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()
	return srv.Shutdown(shutdownCtx)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found (using system environment)")
	}
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
