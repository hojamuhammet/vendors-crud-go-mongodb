package service

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=cinema_service.go -destination=mocks/cinema_service_mock.go

type CinemaService interface {
	GetAllCinemas(page, pageSize int) ([]*domain.GetCinemaResponse, error)
	GetTotalCinemasCount() (int, error)
	GetCinemaByID(id primitive.ObjectID) (*domain.GetCinemaResponse, error)
	CreateCinema(request *domain.CreateCinemaRequest) (*domain.CreateCinemaResponse, error)
	UpdateCinema(id primitive.ObjectID, request *domain.UpdateCinemaRequest) (*domain.UpdateCinemaResponse, error)
	DeleteCinema(id primitive.ObjectID) error
	SearchCinemas(query string, page int, pageSize int) ([]*domain.GetCinemaResponse, error)
	FilterCinemasByTags(tags []string, page int, pageSize int) ([]*domain.GetCinemaResponse, error)
}
