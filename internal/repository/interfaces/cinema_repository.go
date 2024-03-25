package repository

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=cinema_repository.go -destination=mocks/cinema_repository_mock.go

type CinemaRepository interface {
	GetAllCinemas(page, pageSize int) ([]*domain.GetCinemaResponse, error)
	GetTotalCinemasCount() (int, error)
	GetCinemaByID(id primitive.ObjectID) (*domain.GetCinemaResponse, error)
	CreateCinema(request *domain.CreateCinemaRequest) (*domain.CreateCinemaResponse, error)
	UpdateCinema(id primitive.ObjectID, request *domain.UpdateCinemaRequest) (*domain.UpdateCinemaResponse, error)
	DeleteCinema(id primitive.ObjectID) error
	SearchCinemas(query string, page int, pageSize int) ([]*domain.GetCinemaResponse, error)
	FilterCinemasByTags(tags []string, page int, pageSize int) ([]*domain.GetCinemaResponse, error)
}
