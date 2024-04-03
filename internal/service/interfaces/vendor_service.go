package service

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=vendor_service.go -destination=mocks/vendor_service_mock.go

type VendorService interface {
	GetAllVendors(page, pageSize int) ([]*domain.GetVendorResponse, error)
	GetTotalVendorsCount() (int, error)
	GetVendorByID(id primitive.ObjectID) (*domain.GetVendorResponse, error)
	CreateVendor(request *domain.CreateVendorRequest) (*domain.CreateVendorResponse, error)
	UpdateVendor(id primitive.ObjectID, request *domain.UpdateVendorRequest) (*domain.UpdateVendorResponse, error)
	DeleteVendor(id primitive.ObjectID) error
	SearchVendors(query string, page int, pageSize int) ([]*domain.GetVendorResponse, error)
	FilterVendorsByTags(tags []string, page int, pageSize int) ([]*domain.GetVendorResponse, error)
}
