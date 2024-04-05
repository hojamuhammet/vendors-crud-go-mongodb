package repository_test

import (
	"errors"
	"testing"
	"vendors/internal/domain"
	mock_repository "vendors/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetAllVendors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	page := 1
	pageSize := 10

	vendor := &domain.GetVendorResponse{
		ID: primitive.NewObjectID(),
	}

	tests := []struct {
		name  string
		setup func()
		check func([]*domain.GetVendorResponse, error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().GetAllVendors(page, pageSize).Return([]*domain.GetVendorResponse{vendor}, nil)
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, vendor, response[0])
			},
		},
		{
			name: "Find error",
			setup: func() {
				mockVendorRepo.EXPECT().GetAllVendors(page, pageSize).Return(nil, errors.New("find error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "find error", err.Error())
			},
		},
		{
			name: "Decode error",
			setup: func() {
				mockVendorRepo.EXPECT().GetAllVendors(page, pageSize).Return(nil, errors.New("decode error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "decode error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			response, err := mockVendorRepo.GetAllVendors(page, pageSize)
			tt.check(response, err)
		})
	}
}

func TestGetVendorByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	id := primitive.NewObjectID()
	vendor := &domain.GetVendorResponse{
		ID: id,
	}

	tests := []struct {
		name  string
		setup func()
		check func(*domain.GetVendorResponse, error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().GetVendorByID(id).Return(vendor, nil)
			},
			check: func(vendor *domain.GetVendorResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, vendor)
				assert.Equal(t, id, vendor.ID)
			},
		},
		{
			name: "No document found",
			setup: func() {
				mockVendorRepo.EXPECT().GetVendorByID(id).Return(nil, mongo.ErrNoDocuments)
			},
			check: func(vendor *domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, vendor)
				assert.True(t, errors.Is(err, mongo.ErrNoDocuments))
				assert.Equal(t, mongo.ErrNoDocuments.Error(), err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			vendor, err := mockVendorRepo.GetVendorByID(id)
			tt.check(vendor, err)
		})
	}
}

func TestCreateVendor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	vendorRequest := &domain.CreateVendorRequest{
		Cover:          "cover",
		Type:           "type",
		Name:           "name",
		Location:       "location",
		PhoneNumbers:   []string{"1234567890"},
		Websites:       []string{"www.example.com"},
		SocialNetworks: []string{"social_network"},
		Media:          []string{"media"},
		Tags:           []string{"tag"},
		Categories:     []string{"category"},
	}

	vendorResponse := &domain.CreateVendorResponse{
		ID:             primitive.NewObjectID(),
		Cover:          "cover",
		Type:           "type",
		Name:           "name",
		Location:       "location",
		PhoneNumbers:   []string{"1234567890"},
		Websites:       []string{"www.example.com"},
		SocialNetworks: []string{"social_network"},
		Media:          []string{"media"},
		Tags:           []string{"tag"},
		Categories:     []string{"category"},
	}

	tests := []struct {
		name  string
		setup func()
		check func(*domain.CreateVendorResponse, error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().CreateVendor(vendorRequest).Return(vendorResponse, nil)
			},
			check: func(response *domain.CreateVendorResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, vendorResponse, response)
			},
		},
		{
			name: "InsertOne error",
			setup: func() {
				mockVendorRepo.EXPECT().CreateVendor(vendorRequest).Return(nil, errors.New("error inserting vendor document"))
			},
			check: func(response *domain.CreateVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "error inserting vendor document", err.Error())
			},
		},
		{
			name: "InsertedID type assertion error",
			setup: func() {
				mockVendorRepo.EXPECT().CreateVendor(vendorRequest).Return(nil, errors.New("error getting inserted vendor ID"))
			},
			check: func(response *domain.CreateVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "error getting inserted vendor ID", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			response, err := mockVendorRepo.CreateVendor(vendorRequest)
			tt.check(response, err)
		})
	}
}

func TestDeleteVendor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	id := primitive.NewObjectID()

	tests := []struct {
		name  string
		setup func()
		check func(error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().DeleteVendor(id).Return(nil)
			},
			check: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Vendor not found",
			setup: func() {
				mockVendorRepo.EXPECT().DeleteVendor(id).Return(errors.New("vendor not found"))
			},
			check: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, "vendor not found", err.Error())
			},
		},
		{
			name: "DeleteOne error",
			setup: func() {
				mockVendorRepo.EXPECT().DeleteVendor(id).Return(errors.New("delete error"))
			},
			check: func(err error) {
				assert.Error(t, err)
				assert.Equal(t, "delete error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := mockVendorRepo.DeleteVendor(id)
			tt.check(err)
		})
	}
}

func TestSearchVendors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	query := "test"
	page := 1
	pageSize := 10

	vendor := &domain.GetVendorResponse{
		ID: primitive.NewObjectID(),
	}

	tests := []struct {
		name  string
		setup func()
		check func([]*domain.GetVendorResponse, error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().SearchVendors(query, page, pageSize).Return([]*domain.GetVendorResponse{vendor}, nil)
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, vendor, response[0])
			},
		},
		{
			name: "Find error",
			setup: func() {
				mockVendorRepo.EXPECT().SearchVendors(query, page, pageSize).Return(nil, errors.New("find error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "find error", err.Error())
			},
		},
		{
			name: "Decode error",
			setup: func() {
				mockVendorRepo.EXPECT().SearchVendors(query, page, pageSize).Return(nil, errors.New("decode error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "decode error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			response, err := mockVendorRepo.SearchVendors(query, page, pageSize)
			tt.check(response, err)
		})
	}
}

func TestFilterVendorsByTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVendorRepo := mock_repository.NewMockVendorRepository(ctrl)

	tags := []string{"tag1", "tag2"}
	page := 1
	pageSize := 10

	vendor := &domain.GetVendorResponse{
		ID: primitive.NewObjectID(),
	}

	tests := []struct {
		name  string
		setup func()
		check func([]*domain.GetVendorResponse, error)
	}{
		{
			name: "Success",
			setup: func() {
				mockVendorRepo.EXPECT().FilterVendorsByTags(tags, page, pageSize).Return([]*domain.GetVendorResponse{vendor}, nil)
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, vendor, response[0])
			},
		},
		{
			name: "Find error",
			setup: func() {
				mockVendorRepo.EXPECT().FilterVendorsByTags(tags, page, pageSize).Return(nil, errors.New("find error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "find error", err.Error())
			},
		},
		{
			name: "Decode error",
			setup: func() {
				mockVendorRepo.EXPECT().FilterVendorsByTags(tags, page, pageSize).Return(nil, errors.New("decode error"))
			},
			check: func(response []*domain.GetVendorResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
				assert.Equal(t, "decode error", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			response, err := mockVendorRepo.FilterVendorsByTags(tags, page, pageSize)
			tt.check(response, err)
		})
	}
}
