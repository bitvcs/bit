package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nipalab/nipa/db"
	"github.com/nipalab/nipa/internal/config"
	"github.com/nipalab/nipa/internal/http/api"
	"github.com/nipalab/nipa/internal/repository/sqlite"
	"github.com/nipalab/nipa/internal/snow"
	"github.com/nipalab/nipa/internal/usecase"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbConn, err := createDatabaseConnection(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	slog.Info("migrating database...")
	err = db.MigrateUp(dbConn, "sqlite3")
	if err != nil {
		panic(err)
	}
	slog.Info("database migration completed")

	authRepo := sqlite.NewAuthRepository(dbConn)
	userRepo := sqlite.NewUserRepository(dbConn)

	snowUser, err := snow.NewNode(cfg.SnowflakeNodeID)
	if err != nil {
		panic(err)
	}
	reg := &Registry{
		authUsecase: usecase.NewAuth(cfg.JWTKey, userRepo, authRepo),
		userUsecase: usecase.NewUser(snowUser),
	}

	apiApp := api.NewAPI(reg)
	container := apiApp.SetupRoute()

	address := fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpServer := &http.Server{Addr: address, Handler: container}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		<-c
		slog.Info("shutting down server")
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Error("error shutting down server", "error", err)
		}
		os.Exit(0)
	}()

	slog.Info("server is running", "address", address)
	if err := httpServer.Serve(ln); err != http.ErrServerClosed {
		panic(err)
	}
}

func createDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
