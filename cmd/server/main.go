package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/tmf915-api/internal/config"
	"github.com/tmf915-api/internal/db"
	"github.com/tmf915-api/internal/events"
	"github.com/tmf915-api/internal/handlers"
	"github.com/tmf915-api/internal/repository"
	"github.com/tmf915-api/internal/router"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	// Logger
	log := buildLogger(cfg.Log)

	// Database
	database, err := db.New(cfg.Database.DSN, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetime)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer database.Close()
	log.Info().Str("server", "DESKTOP-V2M1TS5").Str("database", "TMF915").Msg("database connected")

	// Repositories
	aiModelSpecRepo := repository.NewAiModelSpecificationRepo(database)
	aiModelRepo := repository.NewAiModelRepo(database)
	aiContractSpecRepo := repository.NewAiContractSpecificationRepo(database)
	aiContractRepo := repository.NewAiContractRepo(database)
	aiContractViolationRepo := repository.NewAiContractViolationRepo(database)
	alarmRepo := repository.NewAlarmRepo(database)
	ruleRepo := repository.NewRuleRepo(database)
	monitorRepo := repository.NewMonitorRepo(database)
	topicRepo := repository.NewTopicRepo(database)
	eventRepo := repository.NewEventRepo(database)
	hubRepo := repository.NewHubRepo(database)
	listenerRepo := repository.NewListenerRepo(database)

	// Event dispatcher (uses Hub table for callbacks)
	dispatcher := events.NewDispatcher(hubRepo, cfg.Events.DispatchTimeout, cfg.Events.MaxRetry, log)

	// Handlers
	aiModelSpecH := handlers.NewAiModelSpecificationHandler(aiModelSpecRepo, dispatcher)
	aiModelH := handlers.NewAiModelHandler(aiModelRepo, dispatcher)
	aiContractSpecH := handlers.NewAiContractSpecificationHandler(aiContractSpecRepo, dispatcher)
	aiContractH := handlers.NewAiContractHandler(aiContractRepo, dispatcher)
	aiContractViolationH := handlers.NewAiContractViolationHandler(aiContractViolationRepo, dispatcher)
	alarmH := handlers.NewAlarmHandler(alarmRepo, dispatcher)
	ruleH := handlers.NewRuleHandler(ruleRepo, dispatcher)
	monitorH := handlers.NewMonitorHandler(monitorRepo, dispatcher)
	topicH := handlers.NewTopicHandler(topicRepo, dispatcher)
	eventH := handlers.NewEventHandler(eventRepo, topicRepo, dispatcher)
	hubH := handlers.NewHubHandler(hubRepo, listenerRepo)

	// Find swagger.json (relative to working directory)
	swaggerPath := "swagger.json"
	if _, err := os.Stat(swaggerPath); errors.Is(err, os.ErrNotExist) {
		log.Warn().Msg("swagger.json not found in working directory — /swagger.json endpoint will 404")
	}

	// Router
	h := router.New(log, cfg.Auth.APIKey, swaggerPath,
		aiModelSpecH, aiModelH, aiContractSpecH, aiContractH, aiContractViolationH,
		alarmH, ruleH, monitorH, topicH, eventH, hubH,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      h,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server
	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("TMF915 API server starting")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
	log.Info().Msg("server stopped")
}

func buildLogger(cfg config.LogConfig) zerolog.Logger {
	var log zerolog.Logger
	if cfg.Format == "console" {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	} else {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
	switch cfg.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	return log
}
