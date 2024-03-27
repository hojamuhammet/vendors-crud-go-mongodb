package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupCafeRouter(cafeRouter *chi.Mux, cafeService *service.CafeService) {
	cafeHandler := handlers.CafeHandler{
		Router:      cafeRouter,
		CafeService: cafeService,
	}

	cafeRouter.Get("/", cafeHandler.GetAllCafesHandler)
	cafeRouter.Get("/{id}", cafeHandler.GetCafeByIDHandler)
	cafeRouter.Post("/", cafeHandler.CreateCafeHandler)
	cafeRouter.Put("/{id}", cafeHandler.UpdateCafeHandler)
	cafeRouter.Delete("/{id}", cafeHandler.DeleteCafe)
	cafeRouter.Get("/search", cafeHandler.SearchCafesHandler)
	cafeRouter.Get("/filter/tags", cafeHandler.FilterCafesByTagsHandler)
}
