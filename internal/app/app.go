package app

import (
	"log/slog"

	grpcapp "github.com/ksenia-samarina/dbKeyValue/internal/app/grpc"
	"github.com/ksenia-samarina/dbKeyValue/internal/services/lsmt/pkg"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	maxMemTableSize uint,
	btreeSize int,
	bloomFilterN uint,
	bloomFilterFp float64,
	storagePath string,
) *App {
	dbService := pkg.NewLSMT(btreeSize, maxMemTableSize, storagePath, bloomFilterN, bloomFilterFp)
	grpcApp := grpcapp.NewApp(log, dbService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
