// Code generated by MockGen. DO NOT EDIT.
// Source: vendor_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"
	domain "vendors/internal/domain"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockVendorRepository is a mock of VendorRepository interface.
type MockVendorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVendorRepositoryMockRecorder
}

// MockVendorRepositoryMockRecorder is the mock recorder for MockVendorRepository.
type MockVendorRepositoryMockRecorder struct {
	mock *MockVendorRepository
}

// NewMockVendorRepository creates a new mock instance.
func NewMockVendorRepository(ctrl *gomock.Controller) *MockVendorRepository {
	mock := &MockVendorRepository{ctrl: ctrl}
	mock.recorder = &MockVendorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVendorRepository) EXPECT() *MockVendorRepositoryMockRecorder {
	return m.recorder
}

// CreateVendor mocks base method.
func (m *MockVendorRepository) CreateVendor(request *domain.CreateVendorRequest) (*domain.CreateVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVendor", request)
	ret0, _ := ret[0].(*domain.CreateVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVendor indicates an expected call of CreateVendor.
func (mr *MockVendorRepositoryMockRecorder) CreateVendor(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVendor", reflect.TypeOf((*MockVendorRepository)(nil).CreateVendor), request)
}

// DeleteVendor mocks base method.
func (m *MockVendorRepository) DeleteVendor(id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVendor", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVendor indicates an expected call of DeleteVendor.
func (mr *MockVendorRepositoryMockRecorder) DeleteVendor(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVendor", reflect.TypeOf((*MockVendorRepository)(nil).DeleteVendor), id)
}

// FilterVendorsByTags mocks base method.
func (m *MockVendorRepository) FilterVendorsByTags(tags []string, page, pageSize int) ([]*domain.GetVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterVendorsByTags", tags, page, pageSize)
	ret0, _ := ret[0].([]*domain.GetVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterVendorsByTags indicates an expected call of FilterVendorsByTags.
func (mr *MockVendorRepositoryMockRecorder) FilterVendorsByTags(tags, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterVendorsByTags", reflect.TypeOf((*MockVendorRepository)(nil).FilterVendorsByTags), tags, page, pageSize)
}

// GetAllVendors mocks base method.
func (m *MockVendorRepository) GetAllVendors(page, pageSize int) ([]*domain.GetVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllVendors", page, pageSize)
	ret0, _ := ret[0].([]*domain.GetVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllVendors indicates an expected call of GetAllVendors.
func (mr *MockVendorRepositoryMockRecorder) GetAllVendors(page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllVendors", reflect.TypeOf((*MockVendorRepository)(nil).GetAllVendors), page, pageSize)
}

// GetTotalVendorsCount mocks base method.
func (m *MockVendorRepository) GetTotalVendorsCount() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalVendorsCount")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalVendorsCount indicates an expected call of GetTotalVendorsCount.
func (mr *MockVendorRepositoryMockRecorder) GetTotalVendorsCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalVendorsCount", reflect.TypeOf((*MockVendorRepository)(nil).GetTotalVendorsCount))
}

// GetVendorByID mocks base method.
func (m *MockVendorRepository) GetVendorByID(id primitive.ObjectID) (*domain.GetVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVendorByID", id)
	ret0, _ := ret[0].(*domain.GetVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVendorByID indicates an expected call of GetVendorByID.
func (mr *MockVendorRepositoryMockRecorder) GetVendorByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVendorByID", reflect.TypeOf((*MockVendorRepository)(nil).GetVendorByID), id)
}

// SearchVendors mocks base method.
func (m *MockVendorRepository) SearchVendors(query string, page, pageSize int) ([]*domain.GetVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchVendors", query, page, pageSize)
	ret0, _ := ret[0].([]*domain.GetVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchVendors indicates an expected call of SearchVendors.
func (mr *MockVendorRepositoryMockRecorder) SearchVendors(query, page, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchVendors", reflect.TypeOf((*MockVendorRepository)(nil).SearchVendors), query, page, pageSize)
}

// UpdateVendor mocks base method.
func (m *MockVendorRepository) UpdateVendor(id primitive.ObjectID, request *domain.UpdateVendorRequest) (*domain.UpdateVendorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVendor", id, request)
	ret0, _ := ret[0].(*domain.UpdateVendorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateVendor indicates an expected call of UpdateVendor.
func (mr *MockVendorRepositoryMockRecorder) UpdateVendor(id, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVendor", reflect.TypeOf((*MockVendorRepository)(nil).UpdateVendor), id, request)
}
