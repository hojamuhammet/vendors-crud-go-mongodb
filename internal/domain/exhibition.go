package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommonExhibitionRequest struct {
	Cover          string   `json:"cover" bson:"cover"`
	Type           string   `json:"type" bson:"type"`
	Name           string   `json:"name" bson:"name"`
	Location       string   `json:"location" bson:"location"`
	PhoneNumbers   []string `json:"phone_numbers" bson:"phone_numbers"`
	Websites       []string `json:"websites" bson:"websites"`
	SocialNetworks []string `json:"social_networks" bson:"social_networks"`
	Media          []string `json:"media" bson:"media"`
	Tags           []string `json:"tags" bson:"tags"`
}

type CommonExhibitionResponse struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Cover          string             `json:"cover" bson:"cover"`
	Type           string             `json:"type" bson:"type"`
	Name           string             `json:"name" bson:"name"`
	Location       string             `json:"location" bson:"location"`
	PhoneNumbers   []string           `json:"phone_numbers" bson:"phone_numbers"`
	Websites       []string           `json:"websites" bson:"websites"`
	SocialNetworks []string           `json:"social_networks" bson:"social_networks"`
	Media          []string           `json:"media" bson:"media"`
	Tags           []string           `json:"tags" bson:"tags"`
}

type GetExhibitionResponse CommonExhibitionResponse
type CreateExhibitionRequest CommonExhibitionRequest
type CreateExhibitionResponse CommonExhibitionResponse
type UpdateExhibitionRequest CommonExhibitionRequest
type UpdateExhibitionResponse CommonExhibitionResponse
