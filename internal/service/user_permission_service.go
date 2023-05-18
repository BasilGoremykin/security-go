package service

import "go-security/internal/model"

type UserPermissionService interface {
	FindByID(id int64) (*model.UserPermissions, error)
	AddPermissionToUser(userId int64, permission *model.Permission) error
	RemovePermissionFromUser(userId int64, permission *model.Permission) error
	DeletePermissions(userId int64) error
}
