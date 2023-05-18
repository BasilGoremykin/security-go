package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisCache(t *testing.T) {
	permissionCache, s := setup(t)
	defer s.Close()

	userId := int64(1)

	initialPermissions := permissionCache.GetPermissions(userId)
	assert.Nil(t, initialPermissions)

	err := permissionCache.SetPermissions(userId, permissions)
	assert.NoError(t, err)

	storedPermissions := permissionCache.GetPermissions(userId)
	assert.Equal(t, permissions, storedPermissions)

	err = permissionCache.DeletePermissions(userId)
	assert.NoError(t, err)

	deletedPermissions := permissionCache.GetPermissions(userId)
	assert.Nil(t, deletedPermissions)
}
