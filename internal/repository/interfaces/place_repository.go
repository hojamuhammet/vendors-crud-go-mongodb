package repository

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=place_repository.go -destination=mocks/place_repository_mock.go

type PlaceRepository interface {
	GetAllPlaces(page, pageSize int) ([]*domain.GetPlaceResponse, error)
	GetTotalPlacesCount() (int, error)
	GetPlaceByID(id primitive.ObjectID) (*domain.GetPlaceResponse, error)
	CreatePlace(request *domain.CreatePlaceRequest) (*domain.CreatePlaceResponse, error)
	UpdatePlace(id primitive.ObjectID, request *domain.UpdatePlaceRequest) (*domain.UpdatePlaceResponse, error)
	DeletePlace(id primitive.ObjectID) error
	SearchPlaces(query string, page int, pageSize int) ([]*domain.GetPlaceResponse, error)
	FilterPlacesByTags(tags []string, page int, pageSize int) ([]*domain.GetPlaceResponse, error)
}
