package model

import (
	"fmt"
	"strings"
)

type UserPermissions struct {
	UserID          int64
	UserPermissions map[Permission]struct{}
}

func (up UserPermissions) String() string {
	var permissions []string
	for permission := range up.UserPermissions {
		permissions = append(permissions, permission.String())
	}
	return fmt.Sprintf("UserPermissions{UserID: %d, Permissions: [%s]}", up.UserID, strings.Join(permissions, ", "))

}
