package db

import (
	"context"
	"github.com/rezaDastrs/banking/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashPassword , err:= util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:     util.RandomOwner(),
		HashPassword: hashPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}

	user , err := testQueries.CreateUser(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashPassword , user.HashPassword)
	require.Equal(t, arg.FullName , user.FullName)
	require.Equal(t, arg.Email , user.Email)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T)  {
	createRandomUser(t)
}

func TestGetUser(t *testing.T)  {
	user1 := createRandomUser(t)
	user2 , err := testQueries.GetUser(context.Background(),user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username , user2.Username)
	require.Equal(t, user1.Email , user2.Email)
	require.Equal(t, user1.HashPassword , user2.HashPassword)
	require.Equal(t, user1.FullName , user2.FullName)
	require.WithinDuration(t, user1.CreatedAt,user2.CreatedAt,time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt,user2.PasswordChangedAt,time.Second)

}