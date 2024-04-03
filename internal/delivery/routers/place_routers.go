package routes

import (
	"vendors/internal/delivery/handlers"
	"vendors/internal/service"

	"github.com/go-chi/chi/v5"
)

func SetupVendorRouter(vendorRouter *chi.Mux, vendorService *service.VendorService) {
	vendorHandler := handlers.VendorHandler{
		Router:        vendorRouter,
		VendorService: vendorService,
	}

	vendorRouter.Get("/", vendorHandler.GetAllVendorsHandler)
	vendorRouter.Get("/{id}", vendorHandler.GetVendorByIDHandler)
	vendorRouter.Post("/", vendorHandler.CreateVendorHandler)
	vendorRouter.Put("/{id}", vendorHandler.UpdateVendorHandler)
	vendorRouter.Delete("/{id}", vendorHandler.DeleteVendor)
	vendorRouter.Get("/search", vendorHandler.SearchVendorsHandler)
	vendorRouter.Get("/filter/tags", vendorHandler.FilterVendorsByTagsHandler)
}
