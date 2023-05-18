package impl

import (
	"errors"
	"fmt"
	"go-security/internal/model"
	"go-security/internal/persistence"
	"go-security/internal/persistence/repository"
	"go-security/internal/persistence/repository/dto"
	"go-security/internal/persistence/repository/repo_err"
	"gorm.io/gorm"
)

type userPermissionPostgresRepository struct {
	db *gorm.DB
}

func NewUserPermissionPostgresRepository(db *gorm.DB) repository.UserPermissionRepository {
	return &userPermissionPostgresRepository{db: db}
}

func (r *userPermissionPostgresRepository) GetAllByID(userId int64) ([]dto.UserPermissionDto, error) {
	var userPermissionsDtos []dto.UserPermissionDto

	if err := r.db.Where("user_id = ?", userId).Find(&userPermissionsDtos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &repo_err.EntityNotFoundError{Msg: fmt.Sprintf("user permissions for user with userId: %d not found in db", userId)}
		}
		return nil, fmt.Errorf("error fetching user permissions from db, %w", err)
	}

	return userPermissionsDtos, nil
}

func (r *userPermissionPostgresRepository) AddPermissionToUser(userId int64, permission *model.Permission) error {
	userPermissionDto := persistence.MapModelToDto(permission, userId)
	if err := r.db.Save(userPermissionDto).Error; err != nil {
		return fmt.Errorf("error %w during adding permissions to user with id: %d", err, userId)
	}
	return nil
}

func (r *userPermissionPostgresRepository) RemovePermissionFromUser(userId int64, permission *model.Permission) error {
	userPermissionDto := persistence.MapModelToDto(permission, userId)
	if err := r.db.Delete(userPermissionDto).Error; err != nil {
		return fmt.Errorf("error %w during removing permission from user with id: %d", err, userId)
	}
	return nil
}

func (r *userPermissionPostgresRepository) DeletePermissions(userId int64) error {
	if err := r.db.Delete(&dto.UserProfileDto{}, "user_id = ?", userId).Error; err != nil {
		return fmt.Errorf("error %w during clearing permissions for user with id: %d", err, userId)
	}
	return nil
}
