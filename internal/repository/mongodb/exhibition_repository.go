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

type MongoDBExhibitionRepository struct {
	collection *mongo.Collection
}

func NewMongoDBExhibitionRepository(collection *mongo.Collection) *MongoDBExhibitionRepository {
	return &MongoDBExhibitionRepository{
		collection: collection,
	}
}

func (r *MongoDBExhibitionRepository) GetAllExhibitions(page, pageSize int) ([]*domain.GetExhibitionResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving exhibition list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var exhibitions []*domain.GetExhibitionResponse
	for cursor.Next(context.Background()) {
		var exhibition domain.GetExhibitionResponse
		if err := cursor.Decode(&exhibition); err != nil {
			slog.Error("Error decoding exhibition: ", utils.Err(err))
			return nil, err
		}
		exhibitions = append(exhibitions, &exhibition)
	}

	return exhibitions, nil
}

func (r *MongoDBExhibitionRepository) GetTotalExhibitionsCount() (int, error) {
	filter := bson.M{}

	totalExhibitions, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total exhibitions count", utils.Err(err))
		return 0, err
	}

	return int(totalExhibitions), nil
}

func (r *MongoDBExhibitionRepository) GetExhibitionByID(id primitive.ObjectID) (*domain.GetExhibitionResponse, error) {
	filter := bson.M{"_id": id}

	var exhibition domain.GetExhibitionResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&exhibition)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting exhibition by ID: %v", utils.Err(err))
		return nil, err
	}
	return &exhibition, nil
}

func (r *MongoDBExhibitionRepository) CreateExhibition(exhibition *domain.CreateExhibitionRequest) (*domain.CreateExhibitionResponse, error) {
	c := domain.CreateExhibitionResponse{
		Cover:          exhibition.Cover,
		Type:           exhibition.Type,
		Name:           exhibition.Name,
		Location:       exhibition.Location,
		PhoneNumbers:   exhibition.PhoneNumbers,
		Websites:       exhibition.Websites,
		SocialNetworks: exhibition.SocialNetworks,
		Media:          exhibition.Media,
		Tags:           exhibition.Tags,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting exhibition document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted exhibition ID")
		return nil, errors.New("error getting inserted exhibition ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBExhibitionRepository) UpdateExhibition(id primitive.ObjectID, update *domain.UpdateExhibitionRequest) (*domain.UpdateExhibitionResponse, error) {
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
		slog.Error("error updating exhibition: ", utils.Err(err))
		return nil, err
	}

	updatedExhibition, err := r.GetExhibitionByID(id)
	if err != nil {
		slog.Error("error fetching updated exhibition: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdateExhibitionResponse{
		ID:             updatedExhibition.ID,
		Cover:          updatedExhibition.Cover,
		Type:           updatedExhibition.Type,
		Name:           updatedExhibition.Name,
		Location:       updatedExhibition.Location,
		PhoneNumbers:   updatedExhibition.PhoneNumbers,
		Websites:       updatedExhibition.Websites,
		SocialNetworks: updatedExhibition.SocialNetworks,
		Media:          updatedExhibition.Media,
		Tags:           updatedExhibition.Tags,
	}

	return updateResponse, nil
}

func (r *MongoDBExhibitionRepository) DeleteExhibition(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting exhibition: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("exhibition not found")
	}

	return nil
}

func (r *MongoDBExhibitionRepository) SearchExhibitions(query string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error) {
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

	var exhibitions []*domain.GetExhibitionResponse

	for cursor.Next(context.Background()) {
		var exhibition domain.GetExhibitionResponse
		if err := cursor.Decode(&exhibition); err != nil {
			return nil, err
		}
		exhibitions = append(exhibitions, &exhibition)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return exhibitions, nil
}

func (r *MongoDBExhibitionRepository) FilterExhibitionsByTags(tags []string, page int, pageSize int) ([]*domain.GetExhibitionResponse, error) {
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

	var exhibitions []*domain.GetExhibitionResponse
	for cursor.Next(context.Background()) {
		var exhibition domain.GetExhibitionResponse
		if err := cursor.Decode(&exhibition); err != nil {
			return nil, err
		}
		exhibitions = append(exhibitions, &exhibition)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return exhibitions, nil
}
