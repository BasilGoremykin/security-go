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

type userProfilePostgresRepository struct {
	db *gorm.DB
}

func NewUserProfilePostgresRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfilePostgresRepository{db: db}
}

func (r *userProfilePostgresRepository) GetUserByMsisdn(msisdn int64) (*dto.UserProfileDto, error) {
	var userProfileDto dto.UserProfileDto

	if err := r.db.Where("msisdn = ?", msisdn).First(&userProfileDto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &repo_err.EntityNotFoundError{Msg: fmt.Sprintf("user with msisdn %d not found in repository", msisdn)}
		}
		return nil, fmt.Errorf("error fetching user from db: %w", err)
	}

	return &userProfileDto, nil
}

func (r *userProfilePostgresRepository) GetUserByGuid(guid int64) (*dto.UserProfileDto, error) {
	var userProfileDto dto.UserProfileDto

	if err := r.db.Where("guid = ?", guid).First(&userProfileDto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &repo_err.EntityNotFoundError{Msg: fmt.Sprintf("user with guid %d not found in repository", guid)}
		}
		return nil, fmt.Errorf("error fetching user from db: %w", err)
	}

	return &userProfileDto, nil
}

func (r *userProfilePostgresRepository) GetUserByUserId(userId int64) (*dto.UserProfileDto, error) {
	var userProfileDto dto.UserProfileDto

	if err := r.db.Where("user_id = ?", userId).First(&userProfileDto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &repo_err.EntityNotFoundError{Msg: fmt.Sprintf("user with userId %d not found in repository", userId)}
		}
		return nil, fmt.Errorf("error fetching user from db: %w", err)
	}

	return &userProfileDto, nil
}

// TODO add kafka signal
func (r *userProfilePostgresRepository) Save(userId int64, ssoUserProfile *model.SsoUserProfile) (*dto.UserProfileDto, error) {
	userProfileDto, err := persistence.MapSsoModelUserToDto(userId, ssoUserProfile)
	if err != nil {
		return nil, err
	}
	if err := r.db.Save(userProfileDto).Error; err != nil {
		return nil, fmt.Errorf("error %w during saving profile for user with guid: %s to db", err, ssoUserProfile.Guid)
	}
	return userProfileDto, nil
}
