package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExhibitionService struct {
	ExhibitionService repository.ExhibitionRepository
}

func NewExhibitionService(exhibitionRepository repository.ExhibitionRepository) *ExhibitionService {
	return &ExhibitionService{ExhibitionService: exhibitionRepository}
}

func (s *ExhibitionService) GetAllExhibitions(page, pageSize int) ([]*domain.GetExhibitionResponse, error) {
	return s.ExhibitionService.GetAllExhibitions(page, pageSize)
}

func (s *ExhibitionService) GetTotalExhibitionsCount() (int, error) {
	return s.ExhibitionService.GetTotalExhibitionsCount()
}

func (s *ExhibitionService) GetExhibitionByID(id primitive.ObjectID) (*domain.GetExhibitionResponse, error) {
	return s.ExhibitionService.GetExhibitionByID(id)
}

func (s *ExhibitionService) CreateExhibition(request *domain.CreateExhibitionRequest) (*domain.CreateExhibitionResponse, error) {
	return s.ExhibitionService.CreateExhibition(request)
}

func (s *ExhibitionService) UpdateExhibition(id primitive.ObjectID, request *domain.UpdateExhibitionRequest) (*domain.UpdateExhibitionResponse, error) {
	return s.ExhibitionService.UpdateExhibition(id, request)
}

func (s *ExhibitionService) DeleteExhibition(id primitive.ObjectID) error {
	return s.ExhibitionService.DeleteExhibition(id)
}

func (s *ExhibitionService) SearchExhibitions(query string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error) {
	return s.ExhibitionService.SearchExhibitions(query, page, pageSize)
}

func (s *ExhibitionService) FilterExhibitionsByTags(tags []string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error) {
	return s.ExhibitionService.FilterExhibitionsByTags(tags, page, pageSize)
}
