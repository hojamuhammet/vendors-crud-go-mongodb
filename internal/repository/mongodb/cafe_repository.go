package repository

import (
	"context"
	"errors"
	"log/slog"
	"vendors/internal/domain"
	"vendors/pkg/lib/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBCafeRepository struct {
	collection *mongo.Collection
}

func NewMongoDBCafeRepository(collection *mongo.Collection) *MongoDBCafeRepository {
	return &MongoDBCafeRepository{
		collection: collection,
	}
}

func (r *MongoDBCafeRepository) GetAllCafes(page, pageSize int) ([]*domain.GetCafeResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving cafe list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var cafes []*domain.GetCafeResponse
	for cursor.Next(context.Background()) {
		var Cafe domain.GetCafeResponse
		if err := cursor.Decode(&Cafe); err != nil {
			slog.Error("Error decoding cafe: ", utils.Err(err))
			return nil, err
		}
		cafes = append(cafes, &Cafe)
	}

	return cafes, nil
}

func (r *MongoDBCafeRepository) GetTotalCafesCount() (int, error) {
	filter := bson.M{}

	totalCafes, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total cafes count", utils.Err(err))
		return 0, err
	}

	return int(totalCafes), nil
}

func (r *MongoDBCafeRepository) GetCafeByID(id primitive.ObjectID) (*domain.GetCafeResponse, error) {
	filter := bson.M{"_id": id}

	var cafe domain.GetCafeResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&cafe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting cafe by ID: %v", utils.Err(err))
		return nil, err
	}
	return &cafe, nil
}

func (r *MongoDBCafeRepository) Createcafe(cafe *domain.CreateCafeRequest) (*domain.CreateCafeResponse, error) {
	c := domain.CreateCafeResponse{
		Cover:          cafe.Cover,
		Type:           cafe.Type,
		Name:           cafe.Name,
		Location:       cafe.Location,
		PhoneNumbers:   cafe.PhoneNumbers,
		Websites:       cafe.Websites,
		SocialNetworks: cafe.SocialNetworks,
		Media:          cafe.Media,
		Tags:           cafe.Tags,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting cafe document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted cafe ID")
		return nil, errors.New("error getting inserted cafe ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBCafeRepository) UpdateCafe(id primitive.ObjectID, update *domain.UpdateCafeRequest) (*domain.UpdateCafeResponse, error) {
	updateFields := bson.M{
		"$set": bson.M{
			"cover":           update.Cover,
			"type":            update.Type,
			"name":            update.Name,
			"location":        update.Location,
			"phone_numbers":   update.PhoneNumbers,
			"websites":        update.Websites,
			"social_networks": update.SocialNetworks,
			"media":           update.Media,
			"tags":            update.Tags,
		},
	}

	filter := bson.M{"_id": id}

	_, err := r.collection.UpdateOne(context.Background(), filter, updateFields)
	if err != nil {
		slog.Error("error updating cafe: ", utils.Err(err))
		return nil, err
	}

	updatedCafe, err := r.GetCafeByID(id)
	if err != nil {
		slog.Error("error fetching updated cafe: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdateCafeResponse{
		ID:             updatedCafe.ID,
		Cover:          updatedCafe.Cover,
		Type:           updatedCafe.Type,
		Name:           updatedCafe.Name,
		Location:       updatedCafe.Location,
		PhoneNumbers:   updatedCafe.PhoneNumbers,
		Websites:       updatedCafe.Websites,
		SocialNetworks: updatedCafe.SocialNetworks,
		Media:          updatedCafe.Media,
		Tags:           updatedCafe.Tags,
	}

	return updateResponse, nil
}

func (r *MongoDBCafeRepository) DeleteCafe(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting cafe: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("cafe not found")
	}

	return nil
}

func (r *MongoDBCafeRepository) SearchCafes(query string, page int, pageSize int) ([]*domain.GetCafeResponse, error) {
	offset := (page - 1) * pageSize

	options := options.Find().SetSkip(int64(offset)).SetLimit(int64(pageSize))

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var cafes []*domain.GetCafeResponse

	for cursor.Next(context.Background()) {
		var cafe domain.GetCafeResponse
		if err := cursor.Decode(&cafe); err != nil {
			return nil, err
		}
		cafes = append(cafes, &cafe)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return cafes, nil
}

func (r *MongoDBCafeRepository) FilterCafesByTags(tags []string, page int, pageSize int) ([]*domain.GetCafeResponse, error) {
	offset := (page - 1) * pageSize

	var tagConditions []bson.M
	for _, tag := range tags {
		tagConditions = append(tagConditions, bson.M{"tags": tag})
	}

	filter := bson.M{"$and": tagConditions}

	options := options.Find().SetSkip(int64(offset)).SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var cafes []*domain.GetCafeResponse
	for cursor.Next(context.Background()) {
		var cafe domain.GetCafeResponse
		if err := cursor.Decode(&cafe); err != nil {
			return nil, err
		}
		cafes = append(cafes, &cafe)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return cafes, nil
}
