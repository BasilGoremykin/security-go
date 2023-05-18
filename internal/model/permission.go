package model

import "fmt"

type Permission struct {
	PermissionId       int8
	PermissionTargetId int64
	TargetType         string
}

func (p Permission) String() string {
	return fmt.Sprintf("{PermissionId : %v, PermissionTargetId: %v, TargetType: %v", p.PermissionId, p.PermissionTargetId, p.TargetType)
}
