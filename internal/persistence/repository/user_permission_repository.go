package repository

import (
	"go-security/internal/model"
	"go-security/internal/persistence/repository/dto"
)

type UserPermissionRepository interface {
	GetAllByID(id int64) ([]dto.UserPermissionDto, error)
	AddPermissionToUser(userId int64, permission *model.Permission) error
	RemovePermissionFromUser(userId int64, permission *model.Permission) error
	DeletePermissions(userId int64) error
}
