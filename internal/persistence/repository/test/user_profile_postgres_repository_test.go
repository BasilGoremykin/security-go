package test

// training tests
import (
	"github.com/stretchr/testify/assert"
	"go-security/internal/persistence/repository/repo_err"

	"testing"
)

func TestUserProfilePostgresRepository_GetUserByGuid(t *testing.T) {
	expectedUserDto := getUserDtoFromSsoModel()

	setupProfileMockQuery(`SELECT * FROM "user_profile" WHERE guid = $1 ORDER BY "user_profile"."id" LIMIT 1`, int64(101))

	userProfileDto, err := profileRepo.GetUserByGuid(int64(101))

	assert.NoError(t, err)
	assert.NotNil(t, userProfileDto)
	assert.Equal(t, expectedUserDto, userProfileDto)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserProfilePostgresRepository_GetUserByMsisdn(t *testing.T) {
	expectedUserDto := getUserDtoFromSsoModel()

	setupProfileMockQuery(`SELECT * FROM "user_profile" WHERE msisdn = $1 ORDER BY "user_profile"."id" LIMIT 1`, int64(79011111111))

	userProfileDto, err := profileRepo.GetUserByMsisdn(int64(79011111111))

	assert.NoError(t, err)
	assert.NotNil(t, userProfileDto)
	assert.Equal(t, expectedUserDto, userProfileDto)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserProfilePostgresRepository_GetUserByUserId(t *testing.T) {
	expectedUserDto := getUserDtoFromSsoModel()

	setupProfileMockQuery(`SELECT * FROM "user_profile" WHERE user_id = $1 ORDER BY "user_profile"."id" LIMIT 1`, int64(1))

	userProfileDto, err := profileRepo.GetUserByUserId(int64(1))

	assert.NoError(t, err)
	assert.NotNil(t, userProfileDto)
	assert.Equal(t, expectedUserDto, userProfileDto)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserProfilePostgresRepository_GetUserByGuidNotFound(t *testing.T) {
	missingGuid := int64(111)

	setupProfileMockQueryNotFound(`SELECT * FROM "user_profile" WHERE guid = $1 ORDER BY "user_profile"."id" LIMIT 1`, missingGuid)

	userProfileDto, err := profileRepo.GetUserByGuid(missingGuid)

	assert.Error(t, err)
	assert.Equal(t, &repo_err.EntityNotFoundError{Msg: "user with guid 111 not found in repository"}, err)
	assert.Nil(t, userProfileDto)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
