package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CafeService struct {
	CafeRepository repository.CafeRepository
}

func NewCafeService(cafeRepository repository.CafeRepository) *CafeService {
	return &CafeService{CafeRepository: cafeRepository}
}

func (s *CafeService) GetAllCafes(page, pageSize int) ([]*domain.GetCafeResponse, error) {
	return s.CafeRepository.GetAllCafes(page, pageSize)
}

func (s *CafeService) GetTotalCafesCount() (int, error) {
	return s.CafeRepository.GetTotalCafesCount()
}

func (s *CafeService) GetCafeByID(id primitive.ObjectID) (*domain.GetCafeResponse, error) {
	return s.CafeRepository.GetCafeByID(id)
}

func (s *CafeService) CreateCafe(request *domain.CreateCafeRequest) (*domain.CreateCafeResponse, error) {
	return s.CafeRepository.CreateCafe(request)
}

func (s *CafeService) UpdateCafe(id primitive.ObjectID, update *domain.UpdateCafeRequest) (*domain.UpdateCafeResponse, error) {
	return s.CafeRepository.UpdateCafe(id, update)
}

func (s *CafeService) DeleteCafe(id primitive.ObjectID) error {
	return s.CafeRepository.DeleteCafe(id)
}

func (s *CafeService) SearchCafes(query string, page int, pageSize int) ([]*domain.GetCafeResponse, error) {
	return s.CafeRepository.SearchCafes(query, page, pageSize)
}

func (s *CafeService) FilterCafesByTags(tags []string, page int, pageSize int) ([]*domain.GetCafeResponse, error) {
	return s.CafeRepository.FilterCafesByTags(tags, page, pageSize)
}
