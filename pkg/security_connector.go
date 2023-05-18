package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	away "go-security/pb/go_security"
	"google.golang.org/grpc"
)

const OWNER_PERMISSION = 1
const INVITE_PERMISSION = 2

type SecurityConnector struct {
	conn   *grpc.ClientConn
	client away.SecurityServiceClient
}

func NewSecurityConnector(properties SecurityConnectorProperties) (*SecurityConnector, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{
		"methodConfig": [{
		  "name": [{"service": "%s"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": %d,
			  "InitialBackoff": "0.5s",
			  "MaxBackoff": "5s",
			  "BackoffMultiplier": 2,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`, "SecurityService", properties.maxRetryAttempts)),
	}

	conn, err := grpc.DialContext(context.Background(), properties.Host+":"+properties.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := away.NewSecurityServiceClient(conn)
	return &SecurityConnector{
		conn:   conn,
		client: client,
	}, nil
}

func (sc *SecurityConnector) Authenticate(token string) (*UserProfile, error) {
	authRequest := &away.AuthRequest{Token: token}

	userProfileBinary, err := sc.client.Authenticate(context.Background(), authRequest)
	if err != nil {
		return nil, fmt.Errorf("error during auth request from server %w", err)
	}

	var userProfile UserProfile
	if err = json.Unmarshal(userProfileBinary.GetResponseBytes(), &userProfile); err != nil {
		return nil, fmt.Errorf("error during deserialization auth response from server")
	}

	return &userProfile, nil
}

func (sc *SecurityConnector) IsOwnerUserId(userId int64, targetId int64, targetType string) (bool, error) {
	userPermissionCheckRequest := &away.UserPermissionCheckRequest{UserId: uint64(userId), PermissionId: OWNER_PERMISSION,
		TargetType: targetType, TargetId: uint64(targetId)}

	binaryResponse, err := sc.client.CheckPermissionForUserId(context.Background(), userPermissionCheckRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling permission check response form server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) IsOwnerToken(token string, targetId int64, targetType string) (bool, error) {
	tokenPermissionCheckRequest := &away.TokenPermissionCheckRequest{Token: token, PermissionId: OWNER_PERMISSION,
		TargetType: targetType, TargetId: uint64(targetId)}

	binaryResponse, err := sc.client.CheckPermissionForToken(context.Background(), tokenPermissionCheckRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling permission check response from server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) CanInviteUserId(userId int64, targetId int64, targetType string) (bool, error) {
	userPermissionCheckRequest := &away.UserPermissionCheckRequest{UserId: uint64(userId), PermissionId: INVITE_PERMISSION,
		TargetType: targetType, TargetId: uint64(targetId)}

	binaryResponse, err := sc.client.CheckPermissionForUserId(context.Background(), userPermissionCheckRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling permission check response form server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) CanInviteToken(token string, targetId int64, targetType string) (bool, error) {
	tokenPermissionCheckRequest := &away.TokenPermissionCheckRequest{Token: token, PermissionId: INVITE_PERMISSION,
		TargetType: targetType, TargetId: uint64(targetId)}

	binaryResponse, err := sc.client.CheckPermissionForToken(context.Background(), tokenPermissionCheckRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling permission check response from server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) AddPermissionToUser(userId int64, permissionId uint32, targetId int64, targetType string) (bool, error) {
	addPermissionRequest := &away.AddPermissionToUserRequest{UserId: uint64(userId), PermissionId: permissionId,
		TargetId: uint64(targetId), TargetType: targetType}

	binaryResponse, err := sc.client.AddPermissionToUser(context.Background(), addPermissionRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling add permission response from server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) RemovePermissionToUser(userId int64, permissionId uint32, targetId int64, targetType string) (bool, error) {
	removePermissionRequest := &away.RemovePermissionFromUserRequest{UserId: uint64(userId), PermissionId: permissionId,
		TargetId: uint64(targetId), TargetType: targetType}

	binaryResponse, err := sc.client.RemovePermissionFromUser(context.Background(), removePermissionRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling remove permission response from server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func (sc *SecurityConnector) ClearUserPermissions(userId int64) (bool, error) {
	clearPermissionsRequest := &away.ClearUsersPermissionsRequest{UserId: uint64(userId)}

	binaryResponse, err := sc.client.ClearUserPermissions(context.Background(), clearPermissionsRequest)
	if err != nil {
		return false, fmt.Errorf("error during handling clear permissions response from server %w", err)
	}

	return parseBoolResponse(binaryResponse)
}

func parseBoolResponse(response *away.SecurityServiceResponse) (bool, error) {
	const (
		trueByte  byte = 1
		falseByte byte = 0
	)

	if len(response.ResponseBytes) == 0 {
		return false, fmt.Errorf("response from security-server is empty")
	}

	switch response.ResponseBytes[0] {
	case trueByte:
		return true, nil
	case falseByte:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected byte in security-server response")
	}
}

func (sc *SecurityConnector) Close() error {
	return sc.conn.Close()
}
