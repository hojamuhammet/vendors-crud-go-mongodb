package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TheatreService struct {
	TheatreService repository.TheatreRepository
}

func NewTheatreService(theatreRepository repository.TheatreRepository) *TheatreService {
	return &TheatreService{TheatreService: theatreRepository}
}

func (s *TheatreService) GetAllTheatres(page, pageSize int) ([]*domain.GetTheatreResponse, error) {
	return s.TheatreService.GetAllTheatres(page, pageSize)
}

func (s *TheatreService) GetTotalTheatresCount() (int, error) {
	return s.TheatreService.GetTotalTheatresCount()
}

func (s *TheatreService) GetTheatreByID(id primitive.ObjectID) (*domain.GetTheatreResponse, error) {
	return s.TheatreService.GetTheatreByID(id)
}

func (s *TheatreService) CreateTheatre(request *domain.CreateTheatreRequest) (*domain.CreateTheatreResponse, error) {
	return s.TheatreService.CreateTheatre(request)
}

func (s *TheatreService) UpdateTheatre(id primitive.ObjectID, request *domain.UpdateTheatreRequest) (*domain.UpdateTheatreResponse, error) {
	return s.TheatreService.UpdateTheatre(id, request)
}

func (s *TheatreService) DeleteTheatre(id primitive.ObjectID) error {
	return s.TheatreService.DeleteTheatre(id)
}

func (s *TheatreService) SearchTheatres(query string, page int, pageSize int) ([]*domain.GetTheatreResponse, error) {
	return s.TheatreService.SearchTheatres(query, page, pageSize)
}

func (s *TheatreService) FilterTheatresByTags(tags []string, page int, pageSize int) ([]*domain.GetTheatreResponse, error) {
	return s.TheatreService.FilterTheatresByTags(tags, page, pageSize)
}
