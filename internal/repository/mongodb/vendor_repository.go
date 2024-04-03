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

type MongoDBVendorRepository struct {
	collection *mongo.Collection
}

func NewMongoDBVendorRepository(collection *mongo.Collection) *MongoDBVendorRepository {
	return &MongoDBVendorRepository{
		collection: collection,
	}
}

func (r *MongoDBVendorRepository) GetAllVendors(page, pageSize int) ([]*domain.GetVendorResponse, error) {
	skip := (page - 1) * pageSize

	filter := bson.M{}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		slog.Error("error retrieving vendors list", utils.Err(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var vendors []*domain.GetVendorResponse
	for cursor.Next(context.Background()) {
		var vendor domain.GetVendorResponse
		if err := cursor.Decode(&vendor); err != nil {
			slog.Error("Error decoding vendor: ", utils.Err(err))
			return nil, err
		}
		vendors = append(vendors, &vendor)
	}

	return vendors, nil
}

func (r *MongoDBVendorRepository) GetTotalVendorsCount() (int, error) {
	filter := bson.M{}

	totalVendors, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		slog.Error("error getting total vendor count", utils.Err(err))
		return 0, err
	}

	return int(totalVendors), nil
}

func (r *MongoDBVendorRepository) GetVendorByID(id primitive.ObjectID) (*domain.GetVendorResponse, error) {
	filter := bson.M{"_id": id}

	var vendor domain.GetVendorResponse

	err := r.collection.FindOne(context.Background(), filter).Decode(&vendor)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		slog.Error("error getting vendor by ID: %v", utils.Err(err))
		return nil, err
	}
	return &vendor, nil
}

func (r *MongoDBVendorRepository) CreateVendor(vendor *domain.CreateVendorRequest) (*domain.CreateVendorResponse, error) {
	c := domain.CreateVendorResponse{
		Cover:          vendor.Cover,
		Type:           vendor.Type,
		Name:           vendor.Name,
		Location:       vendor.Location,
		PhoneNumbers:   vendor.PhoneNumbers,
		Websites:       vendor.Websites,
		SocialNetworks: vendor.SocialNetworks,
		Media:          vendor.Media,
		Tags:           vendor.Tags,
		Categories:     vendor.Categories,
	}

	result, err := r.collection.InsertOne(context.Background(), c)
	if err != nil {
		slog.Error("error inserting vendor document: %v", utils.Err(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		slog.Error("error getting inserted vendor ID")
		return nil, errors.New("error getting inserted vendor ID")
	}

	c.ID = insertedID

	return &c, nil
}

func (r *MongoDBVendorRepository) UpdateVendor(id primitive.ObjectID, update *domain.UpdateVendorRequest) (*domain.UpdateVendorResponse, error) {
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
		slog.Error("error updating vendor: ", utils.Err(err))
		return nil, err
	}

	updatedVendor, err := r.GetVendorByID(id)
	if err != nil {
		slog.Error("error fetching updated vendor: ", utils.Err(err))
		return nil, err
	}

	updateResponse := &domain.UpdateVendorResponse{
		ID:             updatedVendor.ID,
		Cover:          updatedVendor.Cover,
		Type:           updatedVendor.Type,
		Name:           updatedVendor.Name,
		Location:       updatedVendor.Location,
		PhoneNumbers:   updatedVendor.PhoneNumbers,
		Websites:       updatedVendor.Websites,
		SocialNetworks: updatedVendor.SocialNetworks,
		Media:          updatedVendor.Media,
		Tags:           updatedVendor.Tags,
		Categories:     updatedVendor.Categories,
	}

	return updateResponse, nil
}

func (r *MongoDBVendorRepository) DeleteVendor(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		slog.Error("Error deleting vendor: ", utils.Err(err))
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("vendor not found")
	}

	return nil
}

func (r *MongoDBVendorRepository) SearchVendors(query string, page int, pageSize int) ([]*domain.GetVendorResponse, error) {
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

	var vendors []*domain.GetVendorResponse

	for cursor.Next(context.Background()) {
		var vendor domain.GetVendorResponse
		if err := cursor.Decode(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, &vendor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

func (r *MongoDBVendorRepository) FilterVendorsByTags(tags []string, page int, pageSize int) ([]*domain.GetVendorResponse, error) {
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

	var vendors []*domain.GetVendorResponse
	for cursor.Next(context.Background()) {
		var vendor domain.GetVendorResponse
		if err := cursor.Decode(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, &vendor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}
