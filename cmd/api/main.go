package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/di/api"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/postgres"
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadByEnv("transfer")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.Database.DSN, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	producer, err := nsq.NewProducer(cfg.NSQ.NSQDTCPAddr, nsq.NewConfig())
	if err != nil {
		log.Fatalf("failed to create nsq producer: %v", err)
	}
	defer producer.Stop()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Printf("[api] warning: redis unavailable (%v), idempotency will fall back to postgres", err)
	}
	pingCancel()
	defer rdb.Close()

	container := api.BuildContainer(cfg, db, producer, rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", container.Handlers.HTTP.Health)
	mux.HandleFunc("/transfer", container.Handlers.HTTP.CreateTransfer)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("API server listening on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down API server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
