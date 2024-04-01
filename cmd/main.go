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

	log := logger.SetupLogger(cfg.Env)

	slog.Info("Starting the server...", slog.String("env", cfg.Env))
	slog.Debug("Debug messages are enabled") // If env is set to prod, debug messages are going to be disabled

	if err := database.InitDB(cfg); err != nil {
		log.Error("Error setting up MongoDB: %v", err)
	}
	defer database.Close()

	mainRouter := chi.NewRouter()

	cafeRouter := chi.NewRouter()

	mainRouter.Route("/api/cafe", func(r chi.Router) {
		r.Mount("/", cafeRouter)
	})

	cafeCollection := database.GetDB().Collection("cafes")
	cafeRepository := repository.NewMongoDBCafeRepository(cafeCollection)
	cafeService := service.NewCafeService(cafeRepository)
	routes.SetupCafeRouter(cafeRouter, cafeService)

	cinemaRouter := chi.NewRouter()

	mainRouter.Route("/api/cinema", func(r chi.Router) {
		r.Mount("/", cinemaRouter)
	})

	cinemaCollection := database.GetDB().Collection("cinemas")
	cinemaRepository := repository.NewMongoDBCinemaRepository(cinemaCollection)
	cinemaService := service.NewCinemaService(cinemaRepository)
	routes.SetupCinemaRouter(cinemaRouter, cinemaService)

	theatreRouter := chi.NewRouter()

	mainRouter.Route("/api/theatre", func(r chi.Router) {
		r.Mount("/", theatreRouter)
	})

	theatreCollection := database.GetDB().Collection("theatres")
	theatreRepository := repository.NewMongoDBTheatreRepository(theatreCollection)
	theatreService := service.NewTheatreService(theatreRepository)
	routes.SetupTheatreRouter(theatreRouter, theatreService)

	exhibitionRouter := chi.NewRouter()

	mainRouter.Route("/api/exhibition", func(r chi.Router) {
		r.Mount("/", exhibitionRouter)
	})

	exhibitionCollection := database.GetDB().Collection("exhibitions")
	exhibitionRepository := repository.NewMongoDBExhibitionRepository(exhibitionCollection)
	exhibitionService := service.NewExhibitionService(exhibitionRepository)
	routes.SetupExhibitionRouter(exhibitionRouter, exhibitionService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("Shutting down the server gracefully...")

		database.Close()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(cfg.Server.Address, mainRouter); err != nil {
		slog.Error("Server failed to start:", utils.Err(err))
	}
}
