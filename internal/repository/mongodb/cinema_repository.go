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

type MongoDBCinemaRepository struct {
	collection *mongo.Collection
}

func NewMongoDBCinemaRepository(collection *mongo.Collection) *MongoDBCinemaRepository {
	return &MongoDBCinemaRepository{
		collection: collection,
	}
}

func (r *MongoDBCinemaRepository) GetAllCinemas(page, pageSize int) ([]*domain.GetCinemaResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving cinema list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var cinemas []*domain.GetCinemaResponse
	for cursor.Next(context.Background()) {
		var cinema domain.GetCinemaResponse
		if err := cursor.Decode(&cinema); err != nil {
			slog.Error("Error decoding cinema: ", utils.Err(err))
			return nil, err
		}
		cinemas = append(cinemas, &cinema)
	}

	return cinemas, nil
}

func (r *MongoDBCinemaRepository) GetTotalCinemasCount() (int, error) {
	filter := bson.M{}

	totalCinemas, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total cinemas count", utils.Err(err))
		return 0, err
	}

	return int(totalCinemas), nil
}

func (r *MongoDBCinemaRepository) GetCinemaByID(id primitive.ObjectID) (*domain.GetCinemaResponse, error) {
	filter := bson.M{"_id": id}

	var cinema domain.GetCinemaResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&cinema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting cinema by ID: %v", utils.Err(err))
		return nil, err
	}
	return &cinema, nil
}

func (r *MongoDBCinemaRepository) CreateCinema(cinema *domain.CreateCinemaRequest) (*domain.CreateCinemaResponse, error) {
	c := domain.CreateCinemaResponse{
		Cover:          cinema.Cover,
		Type:           cinema.Type,
		Name:           cinema.Name,
		Location:       cinema.Location,
		PhoneNumbers:   cinema.PhoneNumbers,
		Websites:       cinema.Websites,
		SocialNetworks: cinema.SocialNetworks,
		Media:          cinema.Media,
		Tags:           cinema.Tags,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting cinema document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted cinema ID")
		return nil, errors.New("error getting inserted cinema ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBCinemaRepository) UpdateCinema(id primitive.ObjectID, update *domain.UpdateCinemaRequest) (*domain.UpdateCinemaResponse, error) {
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
		slog.Error("error updating cinema: ", utils.Err(err))
		return nil, err
	}

	updatedCinema, err := r.GetCinemaByID(id)
	if err != nil {
		slog.Error("error fetching updated cinema: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdateCinemaResponse{
		ID:             updatedCinema.ID,
		Cover:          updatedCinema.Cover,
		Type:           updatedCinema.Type,
		Name:           updatedCinema.Name,
		Location:       updatedCinema.Location,
		PhoneNumbers:   updatedCinema.PhoneNumbers,
		Websites:       updatedCinema.Websites,
		SocialNetworks: updatedCinema.SocialNetworks,
		Media:          updatedCinema.Media,
		Tags:           updatedCinema.Tags,
	}

	return updateResponse, nil
}

func (r *MongoDBCinemaRepository) DeleteCinema(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting cinema: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("cinema not found")
	}

	return nil
}

func (r *MongoDBCinemaRepository) SearchCinemas(query string, page int, pageSize int) ([]*domain.GetCinemaResponse, error) {
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

	var cinemas []*domain.GetCinemaResponse

	for cursor.Next(context.Background()) {
		var cinema domain.GetCinemaResponse
		if err := cursor.Decode(&cinema); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, &cinema)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return cinemas, nil
}

func (r *MongoDBCinemaRepository) FilterCinemasByTags(tags []string, page int, pageSize int) ([]*domain.GetCinemaResponse, error) {
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

	var cinemas []*domain.GetCinemaResponse
	for cursor.Next(context.Background()) {
		var cinema domain.GetCinemaResponse
		if err := cursor.Decode(&cinema); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, &cinema)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return cinemas, nil
}
