// Code generated by MockGen. DO NOT EDIT.
// Source: authenticate/cmd/data/sqlc (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -destination cmd/data/mock/store.go authenticate/cmd/data/sqlc Store
//

// Package mock_sqlc is a generated GoMock package.
package mock_sqlc

import (
	db "authenticate/cmd/data/sqlc"
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateFavorite mocks base method.
func (m *MockStore) CreateFavorite(arg0 context.Context, arg1 db.CreateFavoriteParams) (db.Favorite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFavorite", arg0, arg1)
	ret0, _ := ret[0].(db.Favorite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFavorite indicates an expected call of CreateFavorite.
func (mr *MockStoreMockRecorder) CreateFavorite(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFavorite", reflect.TypeOf((*MockStore)(nil).CreateFavorite), arg0, arg1)
}

// CreatePlace mocks base method.
func (m *MockStore) CreatePlace(arg0 context.Context, arg1 db.CreatePlaceParams) (db.Place, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlace", arg0, arg1)
	ret0, _ := ret[0].(db.Place)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlace indicates an expected call of CreatePlace.
func (mr *MockStoreMockRecorder) CreatePlace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlace", reflect.TypeOf((*MockStore)(nil).CreatePlace), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// DeletePlace mocks base method.
func (m *MockStore) DeletePlace(arg0 context.Context, arg1 db.DeletePlaceParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlace indicates an expected call of DeletePlace.
func (mr *MockStoreMockRecorder) DeletePlace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlace", reflect.TypeOf((*MockStore)(nil).DeletePlace), arg0, arg1)
}

// DeleteSession mocks base method.
func (m *MockStore) DeleteSession(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockStoreMockRecorder) DeleteSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockStore)(nil).DeleteSession), arg0, arg1)
}

// GetCountryList mocks base method.
func (m *MockStore) GetCountryList(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountryList", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountryList indicates an expected call of GetCountryList.
func (mr *MockStoreMockRecorder) GetCountryList(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountryList", reflect.TypeOf((*MockStore)(nil).GetCountryList), arg0, arg1)
}

// GetFavorite mocks base method.
func (m *MockStore) GetFavorite(arg0 context.Context, arg1 db.GetFavoriteParams) (db.Favorite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorite", arg0, arg1)
	ret0, _ := ret[0].(db.Favorite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorite indicates an expected call of GetFavorite.
func (mr *MockStoreMockRecorder) GetFavorite(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorite", reflect.TypeOf((*MockStore)(nil).GetFavorite), arg0, arg1)
}

// GetPlaceByGoogleId mocks base method.
func (m *MockStore) GetPlaceByGoogleId(arg0 context.Context, arg1 string) (db.Place, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaceByGoogleId", arg0, arg1)
	ret0, _ := ret[0].(db.Place)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaceByGoogleId indicates an expected call of GetPlaceByGoogleId.
func (mr *MockStoreMockRecorder) GetPlaceByGoogleId(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaceByGoogleId", reflect.TypeOf((*MockStore)(nil).GetPlaceByGoogleId), arg0, arg1)
}

// GetRegionList mocks base method.
func (m *MockStore) GetRegionList(arg0 context.Context, arg1 db.GetRegionListParams) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegionList", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegionList indicates an expected call of GetRegionList.
func (mr *MockStoreMockRecorder) GetRegionList(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegionList", reflect.TypeOf((*MockStore)(nil).GetRegionList), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// ListFavoritesByCountrAndRegion mocks base method.
func (m *MockStore) ListFavoritesByCountrAndRegion(arg0 context.Context, arg1 db.ListFavoritesByCountrAndRegionParams) ([]db.ListFavoritesByCountrAndRegionRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFavoritesByCountrAndRegion", arg0, arg1)
	ret0, _ := ret[0].([]db.ListFavoritesByCountrAndRegionRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFavoritesByCountrAndRegion indicates an expected call of ListFavoritesByCountrAndRegion.
func (mr *MockStoreMockRecorder) ListFavoritesByCountrAndRegion(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFavoritesByCountrAndRegion", reflect.TypeOf((*MockStore)(nil).ListFavoritesByCountrAndRegion), arg0, arg1)
}

// ListFavoritesByCountry mocks base method.
func (m *MockStore) ListFavoritesByCountry(arg0 context.Context, arg1 db.ListFavoritesByCountryParams) ([]db.ListFavoritesByCountryRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFavoritesByCountry", arg0, arg1)
	ret0, _ := ret[0].([]db.ListFavoritesByCountryRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFavoritesByCountry indicates an expected call of ListFavoritesByCountry.
func (mr *MockStoreMockRecorder) ListFavoritesByCountry(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFavoritesByCountry", reflect.TypeOf((*MockStore)(nil).ListFavoritesByCountry), arg0, arg1)
}

// ListFavoritesByCreateTime mocks base method.
func (m *MockStore) ListFavoritesByCreateTime(arg0 context.Context, arg1 db.ListFavoritesByCreateTimeParams) ([]db.ListFavoritesByCreateTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFavoritesByCreateTime", arg0, arg1)
	ret0, _ := ret[0].([]db.ListFavoritesByCreateTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFavoritesByCreateTime indicates an expected call of ListFavoritesByCreateTime.
func (mr *MockStoreMockRecorder) ListFavoritesByCreateTime(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFavoritesByCreateTime", reflect.TypeOf((*MockStore)(nil).ListFavoritesByCreateTime), arg0, arg1)
}

// RemoveFavorite mocks base method.
func (m *MockStore) RemoveFavorite(arg0 context.Context, arg1 db.RemoveFavoriteParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFavorite", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFavorite indicates an expected call of RemoveFavorite.
func (mr *MockStoreMockRecorder) RemoveFavorite(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFavorite", reflect.TypeOf((*MockStore)(nil).RemoveFavorite), arg0, arg1)
}

// ToggleFavorite mocks base method.
func (m *MockStore) ToggleFavorite(arg0 context.Context, arg1 db.ToggleFavoriteParams) (db.Favorite, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToggleFavorite", arg0, arg1)
	ret0, _ := ret[0].(db.Favorite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToggleFavorite indicates an expected call of ToggleFavorite.
func (mr *MockStoreMockRecorder) ToggleFavorite(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleFavorite", reflect.TypeOf((*MockStore)(nil).ToggleFavorite), arg0, arg1)
}

// UpdatePlace mocks base method.
func (m *MockStore) UpdatePlace(arg0 context.Context, arg1 db.UpdatePlaceParams) (db.Place, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePlace", arg0, arg1)
	ret0, _ := ret[0].(db.Place)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePlace indicates an expected call of UpdatePlace.
func (mr *MockStoreMockRecorder) UpdatePlace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePlace", reflect.TypeOf((*MockStore)(nil).UpdatePlace), arg0, arg1)
}
