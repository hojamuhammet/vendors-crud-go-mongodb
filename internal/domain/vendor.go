package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommonVendorRequest struct {
	Cover          string   `json:"cover" bson:"cover"`
	Type           string   `json:"type" bson:"type"`
	Name           string   `json:"name" bson:"name"`
	Location       string   `json:"location" bson:"location"`
	PhoneNumbers   []string `json:"phone_numbers" bson:"phone_numbers"`
	Websites       []string `json:"websites" bson:"websites"`
	SocialNetworks []string `json:"social_networks" bson:"social_networks"`
	Media          []string `json:"media" bson:"media"`
	Tags           []string `json:"tags" bson:"tags"`
	Categories     []string `json:"categories" bson:"categories"`
}

type CommonVendorResponse struct {
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
	Categories     []string           `json:"categories" bson:"categories"`
}

type GetVendorResponse CommonVendorResponse
type CreateVendorRequest CommonVendorRequest
type CreateVendorResponse CommonVendorResponse
type UpdateVendorRequest CommonVendorRequest
type UpdateVendorResponse CommonVendorResponse
