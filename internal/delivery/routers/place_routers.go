package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupPlaceRouter(placeRouter *chi.Mux, placeService *service.PlaceService) {
	placeHandler := handlers.PlaceHandler{
		Router:       placeRouter,
		PlaceService: placeService,
	}

	placeRouter.Get("/", placeHandler.GetAllPlacesHandler)
	placeRouter.Get("/{id}", placeHandler.GetPlaceByIDHandler)
	placeRouter.Post("/", placeHandler.CreatePlaceHandler)
	placeRouter.Put("/{id}", placeHandler.UpdatePlaceHandler)
	placeRouter.Delete("/{id}", placeHandler.DeletePlace)
	placeRouter.Get("/search", placeHandler.SearchPlacesHandler)
	placeRouter.Get("/filter/tags", placeHandler.FilterPlacesByTagsHandler)
}
