package dto

type UserPermissionDto struct {
	ID                 int64  `gorm:"column:id"`
	UserID             int64  `gorm:"column:user_id"`
	PermissionID       int8   `gorm:"column:permission_id"`
	PermissionTargetID int64  `gorm:"column:permission_target_id"`
	TargetType         string `gorm:"column:target_type"`
}

func (UserPermissionDto) TableName() string {
	return "user_permission"
}
