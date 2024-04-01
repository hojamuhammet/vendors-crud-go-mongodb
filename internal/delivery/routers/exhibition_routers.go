package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupExhibitionRouter(exhibitionRouter *chi.Mux, exhibitionService *service.ExhibitionService) {
	exhibitionHandler := handlers.ExhibitionHandler{
		Router:            exhibitionRouter,
		ExhibitionService: exhibitionService,
	}

	exhibitionRouter.Get("/", exhibitionHandler.GetAllExhibitionsHandler)
	exhibitionRouter.Get("/{id}", exhibitionHandler.GetExhibitionByIDHandler)
	exhibitionRouter.Post("/", exhibitionHandler.CreateExhibitionHandler)
	exhibitionRouter.Put("/{id}", exhibitionHandler.UpdateExhibitionHandler)
	exhibitionRouter.Delete("/{id}", exhibitionHandler.DeleteExhibition)
	exhibitionRouter.Get("/search", exhibitionHandler.SearchExhibitionsHandler)
	exhibitionRouter.Get("/filter/tags", exhibitionHandler.FilterExhibitionsByTagsHandler)
}
