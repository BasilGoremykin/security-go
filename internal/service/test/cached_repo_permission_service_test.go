package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-security/internal/model"
	"go-security/internal/persistence/repository/dto"
	"go-security/internal/service/impl"
	"testing"
)

type mockUserPermissionRepository struct {
	mock.Mock
}

func (m *mockUserPermissionRepository) GetAllByID(id int64) ([]dto.UserPermissionDto, error) {
	args := m.Called(id)
	return args.Get(0).([]dto.UserPermissionDto), args.Error(1)
}

func (m *mockUserPermissionRepository) AddPermissionToUser(userId int64, permission *model.Permission) error {
	args := m.Called(userId, permission)
	return args.Error(0)
}

func (m *mockUserPermissionRepository) RemovePermissionFromUser(userId int64, permission *model.Permission) error {
	args := m.Called(userId, permission)
	return args.Error(0)
}

func (m *mockUserPermissionRepository) DeletePermissions(userId int64) error {
	args := m.Called(userId)
	return args.Error(0)
}

type mockUserPermissionCache struct {
	mock.Mock
}

func (m *mockUserPermissionCache) GetPermissions(id int64) map[model.Permission]struct{} {
	args := m.Called(id)
	return args.Get(0).(map[model.Permission]struct{})
}

func (m *mockUserPermissionCache) SetPermissions(userId int64, permissions *map[model.Permission]struct{}) error {
	args := m.Called(userId, permissions)
	return args.Error(0)
}

func (m *mockUserPermissionCache) DeletePermissions(userId int64) error {
	args := m.Called(userId)
	return args.Error(0)
}

func TestCacheFound(t *testing.T) {
	repo := new(mockUserPermissionRepository)
	cache := new(mockUserPermissionCache)
	service := impl.NewCachedRepoPermissionService(repo, cache)

	expectedPermissions := map[model.Permission]struct{}{
		model.Permission{PermissionId: 1, PermissionTargetId: 100, TargetType: "location"}: {},
	}

	cache.On("GetPermissions", mock.AnythingOfType("int64")).Return(expectedPermissions)

	userPermissions, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), userPermissions.UserID)
	assert.Equal(t, expectedPermissions, userPermissions.UserPermissions)
	repo.AssertNotCalled(t, "GetAllById")

	cache.AssertExpectations(t)
}

func TestCacheMiss(t *testing.T) {
	repo := new(mockUserPermissionRepository)
	cache := new(mockUserPermissionCache)
	service := impl.NewCachedRepoPermissionService(repo, cache)

	expectedPermissionsDto := []dto.UserPermissionDto{{UserID: 1, PermissionID: 1, PermissionTargetID: 100, TargetType: "location"}}

	emptyMap := make(map[model.Permission]struct{})
	repo.On("GetAllByID", int64(1)).Return(expectedPermissionsDto, nil)
	cache.On("GetPermissions", int64(1)).Return(emptyMap)
	cache.On("SetPermissions", int64(1), mock.Anything).Return(nil)

	userPermissions, err := service.FindByID(int64(1))

	assert.NoError(t, err)
	assert.Equal(t, int64(1), userPermissions.UserID)
	assert.Equal(t, 1, len(userPermissions.UserPermissions))
	cache.AssertCalled(t, "GetPermissions", int64(1))

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestAddPermissionToUser(t *testing.T) {
	repo := new(mockUserPermissionRepository)
	cache := new(mockUserPermissionCache)
	service := impl.NewCachedRepoPermissionService(repo, cache)

	userId := int64(1)
	permission := &model.Permission{
		PermissionId:       1,
		PermissionTargetId: 100,
		TargetType:         "location",
	}

	cache.On("DeletePermissions", userId).Return(nil)
	repo.On("AddPermissionToUser", userId, permission).Return(nil)

	err := service.AddPermissionToUser(userId, permission)

	assert.NoError(t, err)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestDeletePermissions(t *testing.T) {
	repo := new(mockUserPermissionRepository)
	cache := new(mockUserPermissionCache)
	service := impl.NewCachedRepoPermissionService(repo, cache)

	userId := int64(1)

	cache.On("DeletePermissions", userId).Return(nil)
	repo.On("DeletePermissions", userId).Return(nil)

	err := service.DeletePermissions(userId)

	assert.NoError(t, err)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
