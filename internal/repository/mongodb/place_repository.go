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

type MongoDBPlaceRepository struct {
	collection *mongo.Collection
}

func NewMongoDBPlaceRepository(collection *mongo.Collection) *MongoDBPlaceRepository {
	return &MongoDBPlaceRepository{
		collection: collection,
	}
}

func (r *MongoDBPlaceRepository) GetAllPlaces(page, pageSize int) ([]*domain.GetPlaceResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving place list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var places []*domain.GetPlaceResponse
	for cursor.Next(context.Background()) {
		var place domain.GetPlaceResponse
		if err := cursor.Decode(&place); err != nil {
			slog.Error("Error decoding place: ", utils.Err(err))
			return nil, err
		}
		places = append(places, &place)
	}

	return places, nil
}

func (r *MongoDBPlaceRepository) GetTotalPlacesCount() (int, error) {
	filter := bson.M{}

	totalPlaces, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total places count", utils.Err(err))
		return 0, err
	}

	return int(totalPlaces), nil
}

func (r *MongoDBPlaceRepository) GetPlaceByID(id primitive.ObjectID) (*domain.GetPlaceResponse, error) {
	filter := bson.M{"_id": id}

	var place domain.GetPlaceResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&place)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting place by ID: %v", utils.Err(err))
		return nil, err
	}
	return &place, nil
}

func (r *MongoDBPlaceRepository) CreatePlace(place *domain.CreatePlaceRequest) (*domain.CreatePlaceResponse, error) {
	c := domain.CreatePlaceResponse{
		Cover:          place.Cover,
		Type:           place.Type,
		Name:           place.Name,
		Location:       place.Location,
		PhoneNumbers:   place.PhoneNumbers,
		Websites:       place.Websites,
		SocialNetworks: place.SocialNetworks,
		Media:          place.Media,
		Tags:           place.Tags,
		Categories:     place.Categories,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting place document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted place ID")
		return nil, errors.New("error getting inserted place ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBPlaceRepository) UpdatePlace(id primitive.ObjectID, update *domain.UpdatePlaceRequest) (*domain.UpdatePlaceResponse, error) {
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
			"categories":      update.Categories,
		},
	}

	filter := bson.M{"_id": id}

	_, err := r.collection.UpdateOne(context.Background(), filter, updateFields)
	if err != nil {
		slog.Error("error updating place: ", utils.Err(err))
		return nil, err
	}

	updatedPlace, err := r.GetPlaceByID(id)
	if err != nil {
		slog.Error("error fetching updated place: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdatePlaceResponse{
		ID:             updatedPlace.ID,
		Cover:          updatedPlace.Cover,
		Type:           updatedPlace.Type,
		Name:           updatedPlace.Name,
		Location:       updatedPlace.Location,
		PhoneNumbers:   updatedPlace.PhoneNumbers,
		Websites:       updatedPlace.Websites,
		SocialNetworks: updatedPlace.SocialNetworks,
		Media:          updatedPlace.Media,
		Tags:           updatedPlace.Tags,
		Categories:     update.Categories,
	}

	return updateResponse, nil
}

func (r *MongoDBPlaceRepository) DeletePlace(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting place: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("place not found")
	}

	return nil
}

func (r *MongoDBPlaceRepository) SearchPlaces(query string, page int, pageSize int) ([]*domain.GetPlaceResponse, error) {
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

	var places []*domain.GetPlaceResponse

	for cursor.Next(context.Background()) {
		var place domain.GetPlaceResponse
		if err := cursor.Decode(&place); err != nil {
			return nil, err
		}
		places = append(places, &place)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

func (r *MongoDBPlaceRepository) FilterPlacesByTags(tags []string, page int, pageSize int) ([]*domain.GetPlaceResponse, error) {
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

	var places []*domain.GetPlaceResponse
	for cursor.Next(context.Background()) {
		var place domain.GetPlaceResponse
		if err := cursor.Decode(&place); err != nil {
			return nil, err
		}
		places = append(places, &place)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return places, nil
}
