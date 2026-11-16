package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/config"
	"github.com/ferdicsantosa/doitpay-moneytransfer/internal/di/worker"
	"github.com/ferdicsantosa/doitpay-moneytransfer/pkg/postgres"
	"github.com/nsqio/go-nsq"
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

	container := worker.BuildContainer(cfg, db)

	nsqCfg := nsq.NewConfig()
	nsqCfg.MaxAttempts = cfg.NSQ.MaxAttempts
	nsqCfg.MaxInFlight = cfg.NSQ.Concurrency

	consumer, err := nsq.NewConsumer(cfg.NSQ.Topic, cfg.NSQ.Channel, nsqCfg)
	if err != nil {
		log.Fatalf("failed to create nsq consumer: %v", err)
	}
	consumer.AddConcurrentHandlers(container.Handlers.NSQ, cfg.NSQ.Concurrency)

	if err := consumer.ConnectToNSQD(cfg.NSQ.NSQDTCPAddr); err != nil {
		log.Fatalf("failed to connect to nsqd: %v", err)
	}
	log.Printf("Worker listening on topic=%s channel=%s nsqd=%s", cfg.NSQ.Topic, cfg.NSQ.Channel, cfg.NSQ.NSQDTCPAddr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down worker...")
	consumer.Stop()
	<-consumer.StopChan
}
