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

type MongoDBTheatreRepository struct {
	collection *mongo.Collection
}

func NewMongoDBTheatreRepository(collection *mongo.Collection) *MongoDBTheatreRepository {
	return &MongoDBTheatreRepository{
		collection: collection,
	}
}

func (r *MongoDBTheatreRepository) GetAllTheatres(page, pageSize int) ([]*domain.GetTheatreResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving theatre list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var theatres []*domain.GetTheatreResponse
	for cursor.Next(context.Background()) {
		var theatre domain.GetTheatreResponse
		if err := cursor.Decode(&theatre); err != nil {
			slog.Error("Error decoding theatre: ", utils.Err(err))
			return nil, err
		}
		theatres = append(theatres, &theatre)
	}

	return theatres, nil
}

func (r *MongoDBTheatreRepository) GetTotalTheatresCount() (int, error) {
	filter := bson.M{}

	totalTheatres, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total theatres count", utils.Err(err))
		return 0, err
	}

	return int(totalTheatres), nil
}

func (r *MongoDBTheatreRepository) GetTheatreByID(id primitive.ObjectID) (*domain.GetTheatreResponse, error) {
	filter := bson.M{"_id": id}

	var theatre domain.GetTheatreResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&theatre)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting theatre by ID: %v", utils.Err(err))
		return nil, err
	}
	return &theatre, nil
}

func (r *MongoDBTheatreRepository) CreateTheatre(theatre *domain.CreateTheatreRequest) (*domain.CreateTheatreResponse, error) {
	c := domain.CreateTheatreResponse{
		Cover:          theatre.Cover,
		Type:           theatre.Type,
		Name:           theatre.Name,
		Location:       theatre.Location,
		PhoneNumbers:   theatre.PhoneNumbers,
		Websites:       theatre.Websites,
		SocialNetworks: theatre.SocialNetworks,
		Media:          theatre.Media,
		Tags:           theatre.Tags,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting theatre document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted theatre ID")
		return nil, errors.New("error getting inserted theatre ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBTheatreRepository) UpdateTheatre(id primitive.ObjectID, update *domain.UpdateTheatreRequest) (*domain.UpdateTheatreResponse, error) {
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
		slog.Error("error updating theatre: ", utils.Err(err))
		return nil, err
	}

	updatedTheatre, err := r.GetTheatreByID(id)
	if err != nil {
		slog.Error("error fetching updated theatre: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdateTheatreResponse{
		ID:             updatedTheatre.ID,
		Cover:          updatedTheatre.Cover,
		Type:           updatedTheatre.Type,
		Name:           updatedTheatre.Name,
		Location:       updatedTheatre.Location,
		PhoneNumbers:   updatedTheatre.PhoneNumbers,
		Websites:       updatedTheatre.Websites,
		SocialNetworks: updatedTheatre.SocialNetworks,
		Media:          updatedTheatre.Media,
		Tags:           updatedTheatre.Tags,
	}

	return updateResponse, nil
}

func (r *MongoDBTheatreRepository) DeleteTheatre(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting theatre: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("theatre not found")
	}

	return nil
}

func (r *MongoDBTheatreRepository) SearchTheatres(query string, page int, pageSize int) ([]*domain.GetTheatreResponse, error) {
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

	var theatres []*domain.GetTheatreResponse

	for cursor.Next(context.Background()) {
		var theatre domain.GetTheatreResponse
		if err := cursor.Decode(&theatre); err != nil {
			return nil, err
		}
		theatres = append(theatres, &theatre)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return theatres, nil
}

func (r *MongoDBTheatreRepository) FilterTheatresByTags(tags []string, page int, pageSize int) ([]*domain.GetTheatreResponse, error) {
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

	var theatres []*domain.GetTheatreResponse
	for cursor.Next(context.Background()) {
		var theatre domain.GetTheatreResponse
		if err := cursor.Decode(&theatre); err != nil {
			return nil, err
		}
		theatres = append(theatres, &theatre)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return theatres, nil
}
