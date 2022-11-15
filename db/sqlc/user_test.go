package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserPasswordOnly(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, newPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
}

func TestUpdateUser(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, newPassword, newUser.HashedPassword)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, newEmail, newUser.Email)
}
