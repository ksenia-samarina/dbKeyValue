package main

import (
	"github.com/ksenia-samarina/dbKeyValue/internal/app"
	"github.com/ksenia-samarina/dbKeyValue/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	application := app.NewApp(
		log,
		cfg.GRPC.Port,
		cfg.MemTable.MaxSize,
		cfg.MemTable.BtreeSize,
		cfg.BloomFilter.BloomFilterN,
		cfg.BloomFilter.BloomFilterFp,
		cfg.StoragePath,
	)

	go func() {
		application.GRPCServer.MustRun()
	}()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

func setupLogger() *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return log
}
