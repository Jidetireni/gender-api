package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Jidetireni/gender-api/config"
	"github.com/Jidetireni/gender-api/internals/pkg/cache"
	"github.com/Jidetireni/gender-api/internals/pkg/database"
	"github.com/Jidetireni/gender-api/internals/profile"
	"github.com/Jidetireni/gender-api/internals/profile/repository"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error running: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return err
	}

	cache, err := cache.New(cfg)
	if err != nil {
		return err
	}

	profileRepo := repository.NewProfileRepository(db.PostgresDB.DB)
	profileSvc := profile.New(cfg, profileRepo, cache.Redis)

	server := NewServer(
		cfg,
		profileSvc,
		db.PostgresDB,
	)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: server,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil

}
