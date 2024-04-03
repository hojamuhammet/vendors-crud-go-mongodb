package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VendorService struct {
	VendorRepository repository.VendorRepository
}

func NewVendorService(vendorRepository repository.VendorRepository) *VendorService {
	return &VendorService{VendorRepository: vendorRepository}
}

func (s *VendorService) GetAllVendors(page, pageSize int) ([]*domain.GetVendorResponse, error) {
	return s.VendorRepository.GetAllVendors(page, pageSize)
}

func (s *VendorService) GetTotalVendorsCount() (int, error) {
	return s.VendorRepository.GetTotalVendorsCount()
}

func (s *VendorService) GetVendorByID(id primitive.ObjectID) (*domain.GetVendorResponse, error) {
	return s.VendorRepository.GetVendorByID(id)
}

func (s *VendorService) CreateVendor(request *domain.CreateVendorRequest) (*domain.CreateVendorResponse, error) {
	return s.VendorRepository.CreateVendor(request)
}

func (s *VendorService) UpdateVendor(id primitive.ObjectID, update *domain.UpdateVendorRequest) (*domain.UpdateVendorResponse, error) {
	return s.VendorRepository.UpdateVendor(id, update)
}

func (s *VendorService) DeleteVendor(id primitive.ObjectID) error {
	return s.VendorRepository.DeleteVendor(id)
}

func (s *VendorService) SearchVendors(query string, page int, pageSize int) ([]*domain.GetVendorResponse, error) {
	return s.VendorRepository.SearchVendors(query, page, pageSize)
}

func (s *VendorService) FilterVendorsByTags(tags []string, page int, pageSize int) ([]*domain.GetVendorResponse, error) {
	return s.VendorRepository.FilterVendorsByTags(tags, page, pageSize)
}
