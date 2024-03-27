package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupTheatreRouter(theatreRouter *chi.Mux, theatreService *service.TheatreService) {
	theatreHandler := handlers.TheatreHandler{
		Router:         theatreRouter,
		TheatreService: theatreService,
	}

	theatreRouter.Get("/", theatreHandler.GetAllTheatresHandler)
	theatreRouter.Get("/{id}", theatreHandler.GetTheatreByIDHandler)
	theatreRouter.Post("/", theatreHandler.CreateTheatreHandler)
	theatreRouter.Put("/{id}", theatreHandler.UpdateTheatreHandler)
	theatreRouter.Delete("/{id}", theatreHandler.DeleteTheatre)
	theatreRouter.Get("/search", theatreHandler.SearchTheatresHandler)
	theatreRouter.Get("/filter/tags", theatreHandler.FilterTheatresByTagsHandler)
}
