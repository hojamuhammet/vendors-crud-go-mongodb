package service

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=exhibition_service.go -destination=mocks/exhibition_service_mock.go

type ExhibitionService interface {
	GetAllExhibitions(page, pageSize int) ([]*domain.GetExhibitionResponse, error)
	GetTotalExhibitionsCount() (int, error)
	GetExhibitionByID(id primitive.ObjectID) (*domain.GetExhibitionResponse, error)
	CreateExhibition(request *domain.CreateExhibitionRequest) (*domain.CreateExhibitionResponse, error)
	UpdateExhibition(id primitive.ObjectID, request *domain.UpdateExhibitionRequest) (*domain.UpdateExhibitionResponse, error)
	DeleteExhibition(id primitive.ObjectID) error
	SearchExhibitions(query string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error)
	FilterExhibitionsByTags(tags []string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error)
}
