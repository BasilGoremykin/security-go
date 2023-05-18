package cache

import "go-security/internal/model"

type UserPermissionCache interface {
	GetPermissions(userId int64) map[model.Permission]struct{}
	SetPermissions(userId int64, permissions map[model.Permission]struct{}) error
	DeletePermissions(userId int64) error
}
