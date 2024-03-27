package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupCinemaRouter(cinemaRouter *chi.Mux, cinemaService *service.CinemaService) {
	cinemaHandler := handlers.CinemaHandler{
		Router:        cinemaRouter,
		CinemaService: cinemaService,
	}

	cinemaRouter.Get("/", cinemaHandler.GetAllCinemasHandler)
	cinemaRouter.Get("/{id}", cinemaHandler.GetCinemaByIDHandler)
	cinemaRouter.Post("/", cinemaHandler.CreateCinemaHandler)
	cinemaRouter.Put("/{id}", cinemaHandler.UpdateCinemaHandler)
	cinemaRouter.Delete("/{id}", cinemaHandler.DeleteCinema)
	cinemaRouter.Get("/search", cinemaHandler.SearchCinemasHandler)
	cinemaRouter.Get("/filter/tags", cinemaHandler.FilterCinemasByTagsHandler)
}
