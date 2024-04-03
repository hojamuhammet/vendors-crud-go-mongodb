package service

import (
	"vendors/internal/domain"
	repository "vendors/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlaceService struct {
	PlaceRepository repository.PlaceRepository
}

func NewPlaceService(placeRepository repository.PlaceRepository) *PlaceService {
	return &PlaceService{PlaceRepository: placeRepository}
}

func (s *PlaceService) GetAllPlaces(page, pageSize int) ([]*domain.GetPlaceResponse, error) {
	return s.PlaceRepository.GetAllPlaces(page, pageSize)
}

func (s *PlaceService) GetTotalPlacesCount() (int, error) {
	return s.PlaceRepository.GetTotalPlacesCount()
}

func (s *PlaceService) GetPlaceByID(id primitive.ObjectID) (*domain.GetPlaceResponse, error) {
	return s.PlaceRepository.GetPlaceByID(id)
}

func (s *PlaceService) CreatePlace(request *domain.CreatePlaceRequest) (*domain.CreatePlaceResponse, error) {
	return s.PlaceRepository.CreatePlace(request)
}

func (s *PlaceService) UpdatePlace(id primitive.ObjectID, update *domain.UpdatePlaceRequest) (*domain.UpdatePlaceResponse, error) {
	return s.PlaceRepository.UpdatePlace(id, update)
}

func (s *PlaceService) DeletePlace(id primitive.ObjectID) error {
	return s.PlaceRepository.DeletePlace(id)
}

func (s *PlaceService) SearchPlaces(query string, page int, pageSize int) ([]*domain.GetPlaceResponse, error) {
	return s.PlaceRepository.SearchPlaces(query, page, pageSize)
}

func (s *PlaceService) FilterPlacesByTags(tags []string, page int, pageSize int) ([]*domain.GetPlaceResponse, error) {
	return s.PlaceRepository.FilterPlacesByTags(tags, page, pageSize)
}
