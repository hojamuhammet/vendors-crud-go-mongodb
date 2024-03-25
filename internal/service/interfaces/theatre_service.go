package service

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=theatre_service.go -destination=mocks/theatre_service_mock.go

type TheatreService interface {
	GetAllTheatres(page, pageSize int) ([]*domain.GetTheatreResponse, error)
	GetTotalTheatresCount() (int, error)
	GetTheatreByID(id primitive.ObjectID) (*domain.GetTheatreResponse, error)
	CreateTheatre(request *domain.CreateTheatreRequest) (*domain.CreateTheatreResponse, error)
	UpdateTheatre(id primitive.ObjectID, request *domain.UpdateTheatreRequest) (*domain.UpdateTheatreResponse, error)
	DeleteTheatre(id primitive.ObjectID) error
	SearchTheatres(query string, page int, pageSize int) ([]*domain.GetTheatreResponse, error)
	FilterTheatresByTags(tags []string, page int, pageSize int) ([]*domain.GetTheatreResponse, error)
}
