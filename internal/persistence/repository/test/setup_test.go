package test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"go-security/internal/model"
	"go-security/internal/persistence"
	"go-security/internal/persistence/repository"
	"go-security/internal/persistence/repository/dto"
	"go-security/internal/persistence/repository/impl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

var (
	mock           sqlmock.Sqlmock
	profileRepo    repository.UserProfileRepository
	permissionRepo repository.UserPermissionRepository
	ssoUser        = &model.SsoUserProfile{
		IsOrganisation:             "yep",
		Guid:                       "101",
		MnpOperator:                "nil",
		MnpRegion:                  "nil",
		Name:                       "Test",
		GivenName:                  "Test",
		MiddleName:                 "Test",
		LastName:                   "Test",
		Phone:                      "79011111111",
		Email:                      "test@mail.ru",
		AccountNumber:              "1",
		AccountServiceProviderCode: "1",
		Description:                "1",
		Picture:                    "1",
	}
	permission1 = &model.Permission{
		PermissionId:       int8(1),
		PermissionTargetId: int64(100),
		TargetType:         "location",
	}
	permission2 = &model.Permission{
		PermissionId:       int8(2),
		PermissionTargetId: int64(100),
		TargetType:         "location",
	}
	permission3 = &model.Permission{
		PermissionId:       int8(1),
		PermissionTargetId: int64(200),
		TargetType:         "location",
	}
	userWithoutPermissions = int64(3)
)

func setup() {
	var db *sql.DB
	var err error

	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	profileRepo = impl.NewUserProfilePostgresRepository(gormDb)
	permissionRepo = impl.NewUserPermissionPostgresRepository(gormDb)
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setupProfileMockQuery(query string, args ...interface{}) {
	driverValues := make([]driver.Value, len(args))
	for i, v := range args {
		driverValues[i] = v
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(driverValues...).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "guid", "msisdn"}).
			AddRow("1", "101", "79011111111"))
}

func setupProfileMockQueryNotFound(query string, args ...interface{}) {
	driverValues := make([]driver.Value, len(args))
	for i, v := range args {
		driverValues[i] = v
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(driverValues...).
		WillReturnError(gorm.ErrRecordNotFound)
}

func setupGetAllPermissionsBy1UserIdMockQuery() {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_permission" WHERE user_id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "permission_id", "permission_target_id", "target_type"}).
			AddRow("1", "1", "100", "location").
			AddRow("1", "2", "100", "location"))
}

func setupGetAllPermissionsNotFoundMockQuery() {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_permission" WHERE user_id = $1`)).
		WithArgs(userWithoutPermissions).
		WillReturnError(gorm.ErrRecordNotFound)
}

func getUserDtoFromSsoModel() *dto.UserProfileDto {
	expectedUserDto, err := persistence.MapSsoModelUserToDto(int64(1), ssoUser)
	if err != nil {
		panic(err)
	}
	return expectedUserDto
}
