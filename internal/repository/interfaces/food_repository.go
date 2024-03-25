package repository

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=food_repository.go -destination=mocks/food_repository_mock.go

type FoodRepository interface {
	GetAllFoods(page, pageSize int) ([]*domain.GetFoodResponse, error)
	GetTotalFoodsCount() (int, error)
	GetFoodByID(id primitive.ObjectID) (*domain.GetFoodResponse, error)
	CreateFood(request *domain.CreateCinemaRequest) (*domain.CreateFoodResponse, error)
	UpdateFood(id primitive.ObjectID, request *domain.UpdateCinemaRequest) (*domain.UpdateFoodResponse, error)
	DeleteFood(id primitive.ObjectID) error
	SearchFoods(query string, page int, pageSize int) ([]*domain.GetFoodResponse, error)
	FilterFoodsByTags(tags []string, page int, pageSize int) ([]*domain.GetFoodResponse, error)
}
