package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T)  {
	password := RandomString(6)

	hashpassword1 , err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashpassword1)

	err = ChaeckPassword(password , hashpassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = ChaeckPassword(wrongPassword , hashpassword1)
	require.EqualError(t, err , bcrypt.ErrMismatchedHashAndPassword.Error())

	hashpassword2 , err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashpassword2)
	require.NotEqual(t, hashpassword1,hashpassword2)

}
