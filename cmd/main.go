package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vendors/internal/config"
	routes "vendors/internal/delivery/routers"
	repository "vendors/internal/repository/mongodb"
	"vendors/internal/service"
	"vendors/pkg/database"
	"vendors/pkg/lib/utils"
	"vendors/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	logger, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		slog.Error("failed to set up logger: %v", err)
		os.Exit(1)
	}

	slog.Info("Starting the server...", slog.String("env", cfg.Env))
	slog.Debug("Debug messages are enabled")

	if err := database.InitDB(cfg); err != nil {
		logger.ErrorLogger.Error("failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	mainRouter := chi.NewRouter()

	vendorRouter := chi.NewRouter()

	mainRouter.Route("/api/vendor", func(r chi.Router) {
		r.Mount("/", vendorRouter)
	})

	vendorCollection := database.GetDB().Collection("vendors")
	vendorRepository := repository.NewMongoDBVendorRepository(vendorCollection)
	vendorService := service.NewVendorService(vendorRepository)
	routes.SetupVendorRouter(vendorRouter, vendorService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.InfoLogger.Info("Shutting down the server gracefully...")

		database.Close()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(cfg.Server.Address, mainRouter); err != nil {
		logger.ErrorLogger.Error("Server failed to start:", utils.Err(err))
	}
}
