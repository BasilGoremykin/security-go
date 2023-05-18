package dto

type UserProfileDto struct {
	ID     int64  `gorm:"column:id"`
	UserId int64  `gorm:"column:user_id"`
	Guid   string `gorm:"column:guid"`
	Msisdn uint64 `gorm:"column:msisdn"`
}

func (UserProfileDto) TableName() string {
	return "user_profile"
}
