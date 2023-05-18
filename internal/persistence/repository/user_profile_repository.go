package repository

import (
	"go-security/internal/model"
	"go-security/internal/persistence/repository/dto"
)

type UserProfileRepository interface {
	GetUserByGuid(guid int64) (*dto.UserProfileDto, error)
	GetUserByMsisdn(msisdn int64) (*dto.UserProfileDto, error)
	GetUserByUserId(userId int64) (*dto.UserProfileDto, error)
	Save(userId int64, ssoUserProfile *model.SsoUserProfile) (*dto.UserProfileDto, error)
}
