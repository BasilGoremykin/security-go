package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-security/internal/model"
	"go-security/internal/persistence"
	"go-security/internal/persistence/repository"
	"go-security/internal/service"
	away "go-security/pb/go_security"
	"log"
	"strconv"
)

type securityServiceGrpcServer struct {
	away.SecurityServiceServer
	userPermissionsService service.UserPermissionService
	userProfileService     repository.UserProfileRepository
	ssoService             service.OauthService
}

const (
	trueByte  byte = 1
	falseByte byte = 0
)

func returnResponseOnBoolean(b bool) *away.SecurityServiceResponse {
	successResponse := &away.SecurityServiceResponse{ResponseBytes: []byte{trueByte}}
	failureResponse := &away.SecurityServiceResponse{ResponseBytes: []byte{falseByte}}
	if b {
		return successResponse
	}
	return failureResponse
}

func NewSecurityService(permissionService service.UserPermissionService, profileRepository repository.UserProfileRepository, ssoService service.OauthService) away.SecurityServiceServer {
	return &securityServiceGrpcServer{userPermissionsService: permissionService, userProfileService: profileRepository, ssoService: ssoService} // что то тут не так с указателями
}

func (s *securityServiceGrpcServer) Authenticate(ctx context.Context, request *away.AuthRequest) (*away.SecurityServiceResponse, error) {
	token := request.GetToken()

	userId, err := s.getUserIdFromToken(token)
	if err != nil {
		return nil, err
	}

	userProfileDto, err := s.userProfileService.GetUserByUserId(*userId)
	if err != nil {
		ssoProfile, err := s.ssoService.GetSsoProfile(token)
		if err != nil {
			return nil, fmt.Errorf("could not get profile from sso for user %w", err)
		}
		if userProfileDto, err = s.userProfileService.Save(*userId, ssoProfile); err != nil {
			return nil, fmt.Errorf("could not save user profile %w", err)
		}
	}

	userProfileBinary, err := json.Marshal(persistence.MapUserProfileDtoToResponse(userProfileDto))
	if err != nil {
		return nil, fmt.Errorf("could not serialize user profile to binary %w", err)
	}

	return &away.SecurityServiceResponse{ResponseBytes: userProfileBinary}, nil
}

func (s *securityServiceGrpcServer) CheckPermissionForToken(ctx context.Context, request *away.TokenPermissionCheckRequest) (*away.SecurityServiceResponse, error) {
	userId, err := s.getUserIdFromToken(request.GetToken())
	if err != nil {
		return nil, err
	}

	return s.checkPermissionForUser(*userId, int8(request.GetPermissionId()), int64(request.GetTargetId()), request.GetTargetType())
}

func (s *securityServiceGrpcServer) CheckPermissionForUserId(ctx context.Context, request *away.UserPermissionCheckRequest) (*away.SecurityServiceResponse, error) {
	return s.checkPermissionForUser(int64(request.GetUserId()), int8(request.PermissionId), int64(request.GetTargetId()), request.GetTargetType())
}

func (s *securityServiceGrpcServer) AddPermissionToUser(ctx context.Context, request *away.AddPermissionToUserRequest) (*away.SecurityServiceResponse, error) {
	err := s.userPermissionsService.AddPermissionToUser(int64(request.GetUserId()), &model.Permission{PermissionId: int8(request.GetPermissionId()),
		PermissionTargetId: int64(request.GetTargetId()), TargetType: request.GetTargetType()})

	if err != nil {
		logrus.WithError(err).WithField("userId", request.GetUserId()).Warn("Error during adding permission to user")
		return returnResponseOnBoolean(false), nil
	}

	return returnResponseOnBoolean(true), nil
}

func (s *securityServiceGrpcServer) RemovePermissionFromUser(ctx context.Context, request *away.RemovePermissionFromUserRequest) (*away.SecurityServiceResponse, error) {
	err := s.userPermissionsService.RemovePermissionFromUser(int64(request.GetUserId()), &model.Permission{PermissionId: int8(request.GetPermissionId()),
		PermissionTargetId: int64(request.GetTargetId()), TargetType: request.GetTargetType()})

	if err != nil {
		logrus.WithError(err).WithField("userId", request.GetUserId()).Warn("Error during removing permission to user")
		return returnResponseOnBoolean(false), nil
	}

	return returnResponseOnBoolean(true), nil
}

func (s *securityServiceGrpcServer) ClearUserPermissions(ctx context.Context, request *away.ClearUsersPermissionsRequest) (*away.SecurityServiceResponse, error) {
	log.Printf("Trying to clear permissions for user with id: %d", request.GetUserId())

	err := s.userPermissionsService.DeletePermissions(int64(request.GetUserId()))
	if err != nil {
		logrus.WithError(err).WithField("userId", request.GetUserId()).Warn("Error during clearing user's permissions")
		return returnResponseOnBoolean(false), nil
	}

	return returnResponseOnBoolean(true), nil
}

func (s *securityServiceGrpcServer) getUserIdFromToken(token string) (*int64, error) {
	strUserId, err := s.ssoService.Introspect(token)
	if err != nil {
		return nil, fmt.Errorf("could not introspect sso token %w", err)
	}

	userId, err := strconv.ParseInt(strUserId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse sso user_id to int64 %w", err)
	}

	return &userId, nil
}

func (s *securityServiceGrpcServer) checkPermissionForUser(userId int64, permissionId int8, targetId int64, targetType string) (*away.SecurityServiceResponse, error) {
	userPermissions, err := s.userPermissionsService.FindByID(userId)
	if err != nil {
		return nil, err
	}
	permissionToCheck := model.Permission{PermissionId: permissionId,
		PermissionTargetId: targetId,
		TargetType:         targetType}

	_, contains := userPermissions.UserPermissions[permissionToCheck]

	return returnResponseOnBoolean(contains), nil
}
