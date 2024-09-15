package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/b0pof/avito-internship/internal/config"
	delivery "github.com/b0pof/avito-internship/internal/delivery/http"
	"github.com/b0pof/avito-internship/internal/pkg/middleware"
	"github.com/b0pof/avito-internship/internal/repository"
	"github.com/b0pof/avito-internship/internal/server"
	"github.com/b0pof/avito-internship/internal/usecase"
	"github.com/b0pof/avito-internship/pkg/logger"
	"github.com/b0pof/avito-internship/pkg/postgres"
)

const (
	_timeout     = 5 * time.Second
	_connTimeout = 10 * time.Second
)

const (
	_env  = "local"
	_addr = "localhost:8080"
)

type App struct {
	config *config.Config
	server *server.Server
	router *mux.Router
	logger *slog.Logger
}

func MustInit() *App {
	// Config

	cfg := config.MustLoad()

	// issues with SERVER_ADDR variable inside k8s pod
	if cfg.Server.ServerAddr == "" {
		cfg.Server.ServerAddr = _addr
	}

	fmt.Println(cfg)

	// Logger

	log := logger.NewLogger(_env)

	// Router

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok"))
	}).Methods("GET")

	// Postgres

	ctx, cancel := context.WithTimeout(context.Background(), _connTimeout)
	defer cancel()

	pgClient, err := postgres.NewPgxDatabase(ctx, cfg.Postgres)
	if err != nil {
		panic("postgres connection error: " + err.Error())
	}

	// Server

	srv := server.NewServer(cfg.Server, r)

	// Layers

	repo := repository.New(pgClient)
	uc := usecase.New(repo)
	h := delivery.NewHandler(uc)
	h.InitRouter(apiRouter)

	// Middleware
	r.Use(middleware.NewLoggingMiddleware(log))

	return &App{
		config: cfg,
		server: srv,
		router: r,
		logger: log,
	}
}

func (a *App) Run() {
	go func() {
		a.logger.Info("server is running...")
		if err := a.server.Run(); err != nil {
			a.logger.Error("HTTP server ListenAndServe error: " + err.Error())
		}
	}()

	// Graceful shutdown

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	ctx, shutdown := context.WithTimeout(context.Background(), _timeout)
	defer shutdown()

	a.logger.Info("shutting down...")
	if err := a.server.Stop(ctx); err != nil {
		a.logger.Error(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
}
