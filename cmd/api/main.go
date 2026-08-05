package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/seisram/FinApp/internal/application/command"
	"github.com/seisram/FinApp/internal/application/query"
	"github.com/seisram/FinApp/internal/infrastructure/config"
	"github.com/seisram/FinApp/internal/infrastructure/http"
	"github.com/seisram/FinApp/internal/infrastructure/postgres"
)

func main() {
	// Initialize configuration
	cfg := config.Load()

	// Initialize logger
	initLogger(cfg)

	log.Info().Msg("Starting FinApp API server")

	// Initialize database connection pool
	dbpool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database connection pool")
	}
	defer dbpool.Close()

	// Verify database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}
	log.Info().Msg("Database connection established")

	// Initialize repositories
	accountRepo := postgres.NewAccountRepository(dbpool)
	transactionRepo := postgres.NewTransactionRepository(dbpool)
	categoryRepo := postgres.NewCategoryRepository(dbpool)
	partyRepo := postgres.NewPartyRepository(dbpool)

	// Initialize command handlers
	createAccountCmd := command.NewCreateAccountHandler(accountRepo)
	createTransactionCmd := command.NewCreateTransactionHandler(transactionRepo, accountRepo)
	createCategoryCmd := command.NewCreateCategoryHandler(categoryRepo)
	createPartyCmd := command.NewCreatePartyHandler(partyRepo)

	// Initialize query handlers
	accountQuery := query.NewAccountQueryHandler(accountRepo)
	transactionQuery := query.NewTransactionQueryHandler(transactionRepo)
	categoryQuery := query.NewCategoryQueryHandler(categoryRepo)
	partyQuery := query.NewPartyQueryHandler(partyRepo)

	// Initialize HTTP server
	router := gin.Default()
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes
	http.SetupRoutes(router, http.Handlers{
		CreateAccount:      createAccountCmd,
		CreateTransaction:  createTransactionCmd,
		CreateCategory:     createCategoryCmd,
		CreateParty:        createPartyCmd,
		AccountQuery:       accountQuery,
		TransactionQuery:   transactionQuery,
		CategoryQuery:      categoryQuery,
		PartyQuery:         partyQuery,
	})

	// Start server
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	go func() {
		log.Info().Str("address", cfg.Server.Address).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func initLogger(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Log.Format == "json" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}