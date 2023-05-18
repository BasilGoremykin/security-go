package persistence

import (
	"fmt"
	"go-security/internal/model"
	"go-security/internal/persistence/repository/dto"
	"go-security/pkg"
	"strconv"
)

func MapModelToDto(permission *model.Permission, userId int64) *dto.UserPermissionDto {
	return &dto.UserPermissionDto{UserID: userId, PermissionID: permission.PermissionId,
		PermissionTargetID: permission.PermissionTargetId, TargetType: permission.TargetType,
	}
}

func MapDtoToPermission(permissionDto dto.UserPermissionDto) *model.Permission {
	return &model.Permission{PermissionId: permissionDto.PermissionID, PermissionTargetId: permissionDto.PermissionTargetID,
		TargetType: permissionDto.TargetType}
}

func MapDtosToModelSet(dtos []dto.UserPermissionDto) map[model.Permission]struct{} {
	permissionSet := make(map[model.Permission]struct{})
	for _, dto := range dtos {
		permission := MapDtoToPermission(dto)
		permissionSet[*permission] = struct{}{}
	}
	return permissionSet
}

func MapModelSetToArray(set map[model.Permission]struct{}) []model.Permission {
	var permissionsArray []model.Permission
	for permission := range set {
		permissionsArray = append(permissionsArray, permission)
	}
	return permissionsArray
}

func MapModelArrayToModelSet(arr []model.Permission) map[model.Permission]struct{} {
	set := make(map[model.Permission]struct{})
	for _, permission := range arr {
		set[permission] = struct{}{}
	}
	return set
}

func MapSsoModelUserToDto(userId int64, ssoUserProfile *model.SsoUserProfile) (*dto.UserProfileDto, error) {
	msisdn, err := strconv.ParseUint(ssoUserProfile.Phone, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error mapping phone to msisdn %w", err)
	}
	return &dto.UserProfileDto{UserId: userId, Guid: ssoUserProfile.Guid, Msisdn: msisdn}, nil
}

func MapUserProfileDtoToResponse(profileDto *dto.UserProfileDto) *pkg.UserProfile {
	return &pkg.UserProfile{UserId: profileDto.UserId, Guid: profileDto.Guid, Msisdn: profileDto.Msisdn}
}
