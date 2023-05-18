package test

import (
	"github.com/stretchr/testify/assert"
	"go-security/internal/persistence"
	"go-security/internal/persistence/repository/dto"
	"go-security/internal/persistence/repository/repo_err"
	"testing"
)

func TestPermissionPostgresRepository_GetAllById(t *testing.T) {
	expectedPermissionDtos := []dto.UserPermissionDto{*persistence.MapModelToDto(permission1, int64(1)), *persistence.MapModelToDto(permission2, int64(1))}

	setupGetAllPermissionsBy1UserIdMockQuery()

	actualDtos, err := permissionRepo.GetAllByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, actualDtos)
	assert.Equal(t, expectedPermissionDtos, actualDtos)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPermissionPostgresRepository_GetAllByIdNotFound(t *testing.T) {
	setupGetAllPermissionsNotFoundMockQuery()

	nilDtos, err := permissionRepo.GetAllByID(userWithoutPermissions)

	assert.Error(t, err)
	assert.Equal(t, &repo_err.EntityNotFoundError{Msg: "user permissions for user with userId: 3 not found in db"}, err)
	assert.Nil(t, nilDtos)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
