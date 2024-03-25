package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CinemaService struct {
	CinemaRepository repository.CinemaRepository
}

func NewCinemaService(cinemaRepository repository.CinemaRepository) *CinemaService {
	return &CinemaService{CinemaRepository: cinemaRepository}
}

func (s *CinemaService) GetAllCinemas(page, pageSize int) ([]*domain.GetCinemaResponse, error) {
	return s.CinemaRepository.GetAllCinemas(page, pageSize)
}

func (s *CinemaService) GetTotalCinemasCount() (int, error) {
	return s.CinemaRepository.GetTotalCinemasCount()
}

func (s *CinemaService) GetCinemaByID(id primitive.ObjectID) (*domain.GetCinemaResponse, error) {
	return s.CinemaRepository.GetCinemaByID(id)
}

func (s *CinemaService) CreateCinema(request *domain.CreateCinemaRequest) (*domain.CreateCinemaResponse, error) {
	return s.CinemaRepository.CreateCinema(request)
}

func (s *CinemaService) UpdateCinema(id primitive.ObjectID, update *domain.UpdateCinemaRequest) (*domain.UpdateCinemaResponse, error) {
	return s.CinemaRepository.UpdateCinema(id, update)
}

func (s *CinemaService) DeleteCinema(id primitive.ObjectID) error {
	return s.CinemaRepository.DeleteCinema(id)
}

func (s *CinemaService) SearchCinemas(query string, page int, pageSize int) ([]*domain.GetCinemaResponse, error) {
	return s.CinemaRepository.SearchCinemas(query, page, pageSize)
}

func (s *CinemaService) FilterCinemasByTags(tags []string, page int, pageSize int) ([]*domain.GetCinemaResponse, error) {
	return s.CinemaRepository.FilterCinemasByTags(tags, page, pageSize)
}
